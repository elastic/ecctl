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

package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/snaprepoapi"
)

const (
	// floatFormat
	floatFormat = "%.0f"
)

// rpad adds padding to the right of a string if the string is found to be
// empty "-" will be returned instead.
func rpad(s interface{}, padding int) string {
	var str = "-"
	template := fmt.Sprintf("%%-%ds", padding)
	switch t := s.(type) {
	case string:
		if t != "" {
			str = t
		}
	case *string:
		if t != nil {
			str = *t
		}
	case uint16:
		str = strconv.Itoa(int(t))
	case int32:
		str = strconv.Itoa(int(t))
	case *int32:
		str = strconv.Itoa(int(*t))
	case int64:
		str = strconv.Itoa(int(t))
	case int:
		str = strconv.Itoa(t)
	case bool:
		str = strconv.FormatBool(t)
	case *bool:
		if t == nil {
			str = "false"
		} else {
			str = strconv.FormatBool(*t)
		}
	case float64:
		str = strconv.FormatFloat(t, 'f', -1, 64)
	case *float64:
		str = strconv.FormatFloat(*t, 'f', -1, 64)
	default:
		return fmt.Sprintf(template, s)
	}
	return fmt.Sprintf(template, str)
}

// formatBytes formats the capacity to human readable bytes
func formatBytes(rawCap int32, human bool) string {
	const rpadSpace = 5
	if human {
		if rawCap == 0 {
			return "-"
		}
		format := "%.2f"
		if math.Remainder(float64(rawCap), float64(1024)) == 0 {
			format = floatFormat
		}
		capacity := float32(rawCap) / 1024
		if capacity < 1 {
			format = "%d"
			return fmt.Sprint(rpad(fmt.Sprintf(format, rawCap), rpadSpace), "MB")
		}
		if math.Mod(float64(capacity), float64(1024)) == 0 {
			format = floatFormat
		}
		tbCapacity := capacity / 1024
		if tbCapacity < 1 {
			return fmt.Sprint(rpad(fmt.Sprintf(format, capacity), rpadSpace), "GB")
		}
		return fmt.Sprint(rpad(fmt.Sprintf(format, tbCapacity), rpadSpace), "TB")
	}

	return rpad(strconv.Itoa(int(rawCap)), rpadSpace)
}

// formatBytes formats the capacity to human bytes
func formatClusterBytes(rawCap int32, human bool) string {
	var padding = 5
	if human {
		if rawCap == 0 {
			return "-"
		}
		format := "%.2f"
		if math.Remainder(float64(rawCap), float64(1024)) == 0 {
			format = floatFormat
		}
		capacity := float32(rawCap) / 1024
		if capacity < 1 {
			format = "%d"
			return fmt.Sprint(rpad(fmt.Sprintf(format, rawCap), padding), "MB")
		}
		if math.Mod(float64(capacity), float64(1024)) == 0 {
			format = floatFormat
		}
		tbCapacity := capacity / 1024
		if tbCapacity < 1 {
			return fmt.Sprint(rpad(fmt.Sprintf(format, capacity), padding), "GB")
		}
		return fmt.Sprint(rpad(fmt.Sprintf(format, tbCapacity), padding), "TB")
	}

	return rpad(strconv.Itoa(int(rawCap)), 4)
}

func tab() string { return "\t" }

func substr(a, b int32) int32 { return a - b }

func computeClusterCapacity(plan *models.ElasticsearchClusterPlan) int32 {
	var total = int32(0)
	for _, t := range plan.ClusterTopology {
		total += (t.MemoryPerNode * t.NodeCountPerZone) * t.ZoneCount
	}

	return total
}

func derefInt(i *int32) int32 { return *i }
func derefBool(i *bool) bool  { return *i }

func displayAllocator(a *models.AllocatorInfo) bool {
	return *a.Status.Connected || len(a.Instances) > 0
}

func getFailedPlanStepName(plan *models.ElasticsearchClusterPlanInfo) string {
	for _, step := range plan.PlanAttemptLog {
		if *step.Status == "error" && *step.StepID != "plan-completed" {
			return *step.StepID
		}
	}
	return "-"
}

func computePlanDuration(plan *models.ElasticsearchClusterPlanInfo) string {
	start, err := time.Parse(time.RFC3339Nano, plan.AttemptStartTime.String())
	if err != nil {
		return "-"
	}

	end, err := time.Parse(time.RFC3339Nano, plan.AttemptEndTime.String())
	if err != nil {
		return "-"
	}

	return end.Sub(start).String()
}

func getApmFailedPlanStepName(plan *models.ApmPlanInfo) string {
	for _, step := range plan.PlanAttemptLog {
		if *step.Status == "error" && *step.StepID != "plan-completed" {
			return *step.StepID
		}
	}
	return "-"
}

func computeApmPlanDuration(plan *models.ApmPlanInfo) string {
	start, err := time.Parse(time.RFC3339Nano, plan.AttemptStartTime.String())
	if err != nil {
		return "-"
	}

	end, err := time.Parse(time.RFC3339Nano, plan.AttemptEndTime.String())
	if err != nil {
		return "-"
	}

	return end.Sub(start).String()
}

func trimToLen(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n]
}

func rpadTrim(s string, n int) string {
	var pad = n - 2
	if pad < 1 {
		panic("padding is too small, needs to be at least 3")
	}

	return rpad(trimToLen(s, pad), n)
}

// toS3TypeConfig receives an interface and returns an toS3TypeConfig
// the empty values will be transformed to "-"
func toS3TypeConfig(i interface{}) snaprepoapi.S3TypeConfig {
	var buf = new(bytes.Buffer)
	var typeconfig snaprepoapi.S3TypeConfig
	// nolint
	json.NewEncoder(buf).Encode(i)
	// nolint
	json.NewDecoder(buf).Decode(&typeconfig)

	return typeconfig
}

func getClusterName(cluster *models.ElasticsearchClusterInfo) string {
	if cluster.ClusterName == nil || *cluster.ClusterName != *cluster.ClusterID {
		return trimToLen(*cluster.ClusterName, 32)
	}
	return "-"
}

func formatTopologyInfo(clusterInfo models.ElasticsearchClusterInfo) string {
	var format = "%s\t%d"
	var plan = getESCurrentOrPendingPlan(clusterInfo)
	if len(plan.Plan.ClusterTopology) == 0 {
		return fmt.Sprintf(format, "-", 0)
	}

	return fmt.Sprintf(
		format,
		formatClusterBytes(plan.Plan.ClusterTopology[0].MemoryPerNode, true),
		plan.Plan.ClusterTopology[0].ZoneCount,
	)
}

func getESCurrentOrPendingPlan(clusterInfo models.ElasticsearchClusterInfo) *models.ElasticsearchClusterPlanInfo {
	if clusterInfo.PlanInfo.Pending != nil {
		return clusterInfo.PlanInfo.Pending
	}
	// This will be returned if both (Current and Pending) are nil
	// It populates the info in a best-effort manner, so data might
	// be slightly off.
	if clusterInfo.PlanInfo.Current == nil {
		var (
			zoneCount        int32
			zoneNames        []string
			memoryPerNode    int32
			version          = "?.?.?"
			nodeCountPerZone = int32(0)
		)
		for _, i := range clusterInfo.Topology.Instances {
			if !stringInSlice(i.Zone, zoneNames) {
				zoneNames = append(zoneNames, i.Zone)
				zoneCount++
				if strings.Split(*i.InstanceName, "-")[0] == "instance" && memoryPerNode == 0 {
					if i.Memory != nil {
						memoryPerNode = *i.Memory.InstanceCapacity
					}
				}
			}
		}

		if len(clusterInfo.Topology.Instances) > 0 || zoneCount > 0 {
			nodeCountPerZone = int32(len(clusterInfo.Topology.Instances)) / zoneCount
		}

		if len(clusterInfo.Topology.Instances) > 0 {
			version = clusterInfo.Topology.Instances[0].ServiceVersion
		}

		return &models.ElasticsearchClusterPlanInfo{
			Plan: &models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: version,
				},
				ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
					{
						NodeCountPerZone: nodeCountPerZone,
						MemoryPerNode:    memoryPerNode,
						ZoneCount: zoneCount,
					},
				},
			},
		}
	}
	return clusterInfo.PlanInfo.Current
}

func centiCentsToCents(i int) int {
	return i / 100
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// equal checks if the passed values are equal. Currently only strings
// are supported.
func equal(x, y interface{}) bool {
	// setting a to something different than "" to avoid a case where
	// the passed types are not handled by the switch cases
	var a = "x"
	var b string
	switch s := x.(type) {
	case string:
		a = s
	case *string:
		a = *s
	}

	switch s := y.(type) {
	case string:
		b = s
	case *string:
		b = *s
	}

	return a == b
}
