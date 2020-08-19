// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	docsLocation        string
	completionsLocation string
	generatedBinary     string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates completions and docs",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var completionCmd = &cobra.Command{
	Use:     "completions",
	Short:   "Outputs the Bash completion to either stdout (default) or to a file",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		output := os.Stdout
		if completionsLocation != "" {
			var err error
			if output, err = os.Create(completionsLocation); err != nil {
				return err
			}
			defer output.Close()
		}

		RootCmd.Use = generatedBinary
		if strings.Contains(os.Getenv("SHELL"), "zsh") {
			return generateZshCompletion(output, generatedBinary)
		}

		return RootCmd.GenBashCompletion(output)
	},
}

var docCmd = &cobra.Command{
	Use:     "docs",
	Short:   "Generates the command tree documentation",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(docsLocation); os.IsNotExist(err) {
			if err := os.MkdirAll(docsLocation, 0756); err != nil {
				return err
			}
		}

		return genMarkdownTree(RootCmd, docsLocation)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(completionCmd)
	generateCmd.AddCommand(docCmd)
	completionCmd.Flags().StringVarP(&completionsLocation, "location", "l", "", "Sets the location of the generated output")
	completionCmd.Flags().StringVar(&generatedBinary, "binary", "ecctl", "Binary name to set for the autocompletion")
	docCmd.Flags().StringVarP(&docsLocation, "location", "l", "./docs", "Sets the location of the generated output")
}

// generateZshCompletion is a modified version of kubectl completion cmd:
// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubectl/cmd/completion.go#L141
// All credit goes to the the contributors of that file.
func generateZshCompletion(out io.Writer, binary string) error {
	zshInitialization := strings.ReplaceAll(`
__BINARY_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__BINARY_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__BINARY_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__BINARY_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__BINARY_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__BINARY_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__BINARY_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__BINARY_filedir() {
	local RET OLD_IFS w qw

	__BINARY_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __BINARY_debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__BINARY_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__BINARY_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__BINARY_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__BINARY_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__BINARY_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__BINARY_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__BINARY_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__BINARY_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__BINARY_type/g" \
	<<'BASH_COMPLETION_EOF'
`, "BINARY", binary)

	if _, err := out.Write([]byte(zshInitialization)); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := RootCmd.GenBashCompletion(buf); err != nil {
		return err
	}

	if _, err := out.Write(buf.Bytes()); err != nil {
		return err
	}

	zshTail := strings.ReplaceAll(`
BASH_COMPLETION_EOF
}
__BINARY_bash_source <(__BINARY_convert_bash_to_zsh)
`, "BINARY", binary)
	if _, err := out.Write([]byte(zshTail)); err != nil {
		return err
	}

	return nil
}
