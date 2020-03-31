// Copyright 2019 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
package chubaofs

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/utils/mount"
	"os"
)

type nodeServer struct {
	UnimplementedNodeServer
	nodeID string
	driver *driver
}

func NewNodeServer(driver *driver) *nodeServer {
	return &nodeServer{
		driver: driver,
		nodeID: driver.nodeID,
	}
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	hasMount, err := HasMount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if hasMount {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	volumeId := req.GetVolumeId()
	param := req.GetVolumeAttributes()
	cfsServer, err := newCfsServer(volumeId, param)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = cfsServer.persistClientConf(targetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "persist client config file fail, err: %v", err)
	}

	if err = doMount(cfsServer); err != nil {
		return nil, status.Errorf(codes.Internal, "mount fail, err: %v", err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("targetPath[%v] not found", targetPath))
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	//assuming success if already unmounted
	if notMnt {
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}

	err = mount.New("").Unmount(req.GetTargetPath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func doMount(cfsServer *cfsServer) error {
	return cfsServer.runClient()
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.nodeID,
	}, nil
}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}
