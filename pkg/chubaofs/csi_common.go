package chubaofs

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnimplementedControllerServer can be embedded to have forward compatible implementations.
type UnimplementedControllerServer struct {
}

func (*UnimplementedControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVolume not implemented")
}

func (*UnimplementedControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVolume not implemented")
}
func (*UnimplementedControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ControllerPublishVolume not implemented")
}
func (*UnimplementedControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ControllerUnpublishVolume not implemented")
}
func (*UnimplementedControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateVolumeCapabilities not implemented")
}
func (*UnimplementedControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVolumes not implemented")
}
func (*UnimplementedControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCapacity not implemented")
}
func (*UnimplementedControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ControllerGetCapabilities not implemented")
}
func (*UnimplementedControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSnapshot not implemented")
}
func (*UnimplementedControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSnapshot not implemented")
}
func (*UnimplementedControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSnapshots not implemented")
}

// UnimplementedIdentityServer can be embedded to have forward compatible implementations.
type UnimplementedIdentityServer struct {
}

func (*UnimplementedIdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginInfo not implemented")
}
func (*UnimplementedIdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginCapabilities not implemented")
}
func (*UnimplementedIdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Probe not implemented")
}

// UnimplementedNodeServer can be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (ns *UnimplementedNodeServer) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeStageVolume not implemented")
}
func (ns *UnimplementedNodeServer) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeUnstageVolume not implemented")
}

func (ns *UnimplementedNodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodePublishVolume not implemented")
}

func (ns *UnimplementedNodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeUnpublishVolume not implemented")
}

func (ns *UnimplementedNodeServer) NodeGetId(ctx context.Context, req *csi.NodeGetIdRequest) (*csi.NodeGetIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeGetId not implemented")
}

func (ns *UnimplementedNodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeGetInfo not implemented")
}

func (ns *UnimplementedNodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeGetCapabilities not implemented")
}
