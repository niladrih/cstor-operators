// Copyright 2020 The OpenEBS Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"fmt"

	"k8s.io/klog"

	"github.com/openebs/api/v3/pkg/proto"
	"github.com/openebs/api/v3/pkg/util"
	"golang.org/x/net/context"
)

// sudo $ISTGTCONTROL snapdestroy vol1 snapname1 0
// sudo $ISTGTCONTROL snapcreate vol1 snapname1

// constants
const (
	VolumeGrpcListenPort = 7777
	CmdSnapCreate        = "SNAPCREATE"
	CmdSnapDestroy       = "SNAPDESTROY"
	//IoWaitTime is the time interval for which the IO has to be stopped before doing snapshot operation
	IoWaitTime = 10
	//TotalWaitTime is the max time duration to wait for doing snapshot operation on all the replicas
	TotalWaitTime   = 60
	ProtocolVersion = 1
)

// CommandStatus is the response from istgt for control commands
type CommandStatus struct {
	Response string `json:"response"`
}

// APIUnixSockVar is unix socker variable
var APIUnixSockVar util.UnixSock

// Server represents the gRPC server
type Server struct {
}

func init() {
	APIUnixSockVar = util.RealUnixSock{}
}

// RunVolumeSnapCreateCommand performs snapshot create operation and sends back the response
func (s *Server) RunVolumeSnapCreateCommand(ctx context.Context, in *proto.VolumeSnapCreateRequest) (*proto.VolumeSnapCreateResponse, error) {
	klog.Infof("Received snapshot create request. volname = %s, snapname = %s, version = %d", in.Volume, in.Snapname, in.Version)
	volcmd, err := CreateSnapshot(ctx, in)
	return volcmd, err

}

// RunVolumeSnapDeleteCommand performs snapshot create operation and sends back the response
func (s *Server) RunVolumeSnapDeleteCommand(ctx context.Context, in *proto.VolumeSnapDeleteRequest) (*proto.VolumeSnapDeleteResponse, error) {
	klog.Infof("Received snapshot delete request. volname = %s, snapname = %s, version = %d", in.Volume, in.Snapname, in.Version)
	volcmd, err := DeleteSnapshot(ctx, in)
	return volcmd, err
}

// CreateSnapshot sends snapcreate command to istgt and returns the response
func CreateSnapshot(ctx context.Context, in *proto.VolumeSnapCreateRequest) (*proto.VolumeSnapCreateResponse, error) {
	sockresp, err := APIUnixSockVar.SendCommand(fmt.Sprintf("%s %s %s %v %v",
		CmdSnapCreate, in.Volume, in.Snapname, IoWaitTime, TotalWaitTime))
	respstr := "ERR"
	if nil != sockresp && len(sockresp) > 1 {
		respstr = sockresp[1]
	}
	status := CommandStatus{
		Response: respstr,
	}
	jsonresp, _ := json.Marshal(status)
	resp := &proto.VolumeSnapCreateResponse{
		Status: jsonresp,
	}
	return resp, err
}

// DeleteSnapshot sends snapdelete command to istgt and returns the response
func DeleteSnapshot(ctx context.Context, in *proto.VolumeSnapDeleteRequest) (*proto.VolumeSnapDeleteResponse, error) {
	sockresp, err := APIUnixSockVar.SendCommand(fmt.Sprintf("%s %s %s %v %v",
		CmdSnapDestroy, in.Volume, in.Snapname, IoWaitTime, TotalWaitTime))
	respstr := "ERR"
	if nil != sockresp && len(sockresp) > 1 {
		respstr = sockresp[1]
	}
	status := CommandStatus{
		Response: respstr,
	}
	jsonresp, _ := json.Marshal(status)
	resp := &proto.VolumeSnapDeleteResponse{
		Status: jsonresp,
	}
	return resp, err
}
