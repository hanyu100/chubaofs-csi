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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/utils/mount"
)

type nodeServer struct {
	*csi.UnimplementedNodeServer
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
	klog.Infof("NodePublishVolume req:%v", req)
	targetPath := req.GetTargetPath()
	hasMount, err := hasMount(targetPath)
	if err != nil {
		klog.Errorf("check mount status error, %v", err)
		return nil, status.Errorf(codes.Internal, "check mount status error, %v", err)
	}

	if hasMount {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	volumeId := req.GetVolumeId()
	param := req.GetVolumeContext()
	cfsServer, err := newCfsServer(volumeId, param)
	if err != nil {
		klog.Errorf("new cfs server error, %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "new cfs server error, %v", err)
	}

	err = cfsServer.persistClientConf(req.GetTargetPath())
	if err != nil {
		klog.Errorf("persist client config file fail, err: %v", err)
		return nil, status.Errorf(codes.Internal, "persist client config file fail, err: %v", err)
	}

	if err = doMount(cfsServer); err != nil {
		klog.Errorf("mount fail, err: %v", err)
		return nil, status.Errorf(codes.Internal, "mount fail, err: %v", err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.Infof("NodeUnpublishVolume req:%v", req)
	targetPath := req.GetTargetPath()
	hasMount, err := hasMount(targetPath)
	if err != nil {
		klog.Errorf("check mount status error, %v", err)
		return nil, status.Errorf(codes.Internal, "check mount status error, %v", err)
	}

	if hasMount {
		err = mount.New("").Unmount(req.GetTargetPath())
		if err != nil {
			klog.Errorf("unmount error, %v", err)
			return nil, status.Errorf(codes.Internal, "unmount error, %v", err)
		}
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func doMount(cfsServer *cfsServer) error {
	return cfsServer.runClient()
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.Infof("NodeStageVolume req:%v", req)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.Infof("NodeUnstageVolume req:%v", req)
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.Infof("NodeGetInfo req:%v", req)
	return &csi.NodeGetInfoResponse{
		NodeId: ns.nodeID,
	}, nil
}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.Infof("NodeGetCapabilities req:%v", req)
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
