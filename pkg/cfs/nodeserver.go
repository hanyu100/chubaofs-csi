/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cfs

import (
	"encoding/json"
	"github.com/chubaofs/chubaofs-csi/pkg/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"k8s.io/klog"
	"os"
	"path"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
	masterAddress string
}

func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	stagingTargetPath := req.GetStagingTargetPath()
	targetPath := req.GetTargetPath()
	if err := createMountPoint(targetPath); err != nil {
		klog.Errorf("failed to create mount point at %s: %v", targetPath, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	hasMount, err := isMountPoint(targetPath)
	if err != nil {
		klog.Errorf("check mount status error, %v", err)
		return nil, status.Errorf(codes.Internal, "check mount status error, %v", err)
	}

	if hasMount {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	err = bindMount(stagingTargetPath, targetPath)
	if err != nil {
		klog.Errorf("mount -bind stagingTargetPath[%v] to targetPath[%v], %v", stagingTargetPath, targetPath, err)
		return nil, status.Errorf(codes.Internal, "check mount status error, %v", err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	err := umountVolume(targetPath)
	if err != nil {
		klog.Errorf("umount targetPath[%v] fail, %v", targetPath, err)
		return nil, status.Errorf(codes.Internal, "umount targetPath[%v] fail, %v", targetPath, err)
	}

	_ = os.Remove(targetPath)
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	stagingTargetPath := req.GetStagingTargetPath()
	volumeName := req.GetVolumeId()
	err := createMountPoint(stagingTargetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create stagingTargetPath[%v] fail, %v", stagingTargetPath, err)
	}

	hasMount, err := isMountPoint(stagingTargetPath)
	if err != nil {
		klog.Errorf("check mount status error, %v", err)
		return nil, status.Errorf(codes.Internal, "check mount status error, %v", err)
	}

	if hasMount {
		return &csi.NodeStageVolumeResponse{}, nil
	}

	cfgmap := make(map[string]interface{})
	cfgmap[FUSE_KEY_MOUNT_POINT] = stagingTargetPath
	cfgmap[FUSE_KEY_VOLUME_NAME] = volumeName
	cfgmap[FUSE_KEY_MASTER_ADDR] = ns.masterAddress
	cfgmap[FUSE_KEY_LOG_PATH] = "/cfs/conf/" + volumeName
	cfgmap[FUSE_KEY_LOG_LEVEL] = "error"
	cfgmap[FUSE_KEY_LOOKUP_VALID] = "30"
	cfgmap[FUSE_KEY_OWNER] = "cfs"
	cfgmap[FUSE_KEY_PROF_PORT] = "10094"
	//the parameters below are all set by default value
	cfgmap[FUSE_KEY_ICACHE_TIMEOUT] = ""
	cfgmap[FUSE_KEY_ATTR_VALID] = ""
	cfgmap[FUSE_KEY_EN_SYNC_WRITE] = ""
	cfgmap[FUSE_KEY_AUTO_INVAL_DATA] = ""
	cfgmap[FUSE_KEY_RDONLY] = "false"
	cfgmap[FUSE_KEY_WRITE_CACHE] = "false"
	cfgmap[FUSE_KEY_KEEP_CACHE] = "false"

	configFilePath := "/cfs/conf/" + volumeName
	cfgstr, err := json.MarshalIndent(cfgmap, "", "      ")
	err = ioutil.WriteFile(configFilePath, cfgstr, 0444)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create client config file fail. err: %v", err)
	}

	klog.V(0).Infof("create client config file success, volumeId:%v", volumeName)
	if err = mountVolume(configFilePath); err != nil {
		klog.Errorf("mount fail, err: %v", err)
		return nil, status.Errorf(codes.Internal, "mount fail, err: %v", err)
	}

	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	stagingTargetPath := req.GetStagingTargetPath()
	err := umountVolume(stagingTargetPath)
	if err != nil {
		klog.Errorf("umount stagingTargetPath[%v] fail, %v", stagingTargetPath, err)
		return nil, status.Errorf(codes.Internal, "umount stagingTargetPath[%v] fail, %v", stagingTargetPath, err)
	}

	_ = os.Remove(stagingTargetPath)
	return &csi.NodeUnstageVolumeResponse{}, nil
}
