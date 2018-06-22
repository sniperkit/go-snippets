package main

import (
	"image"
)

type Pose struct {
}

type PoseEstimator struct {
}

func NewPoseEstimator() *PoseEstimator {
	return new(PoseEstimator)
}

func DefaultPoseEstimator() *PoseEstimator {
	return new(PoseEstimator)
}

func (PoseEstimator) Estimate(im image.Image) Pose {
	var pose Pose
	return pose
}

func main() {
	// im = imread('myphoto.jpg');
	// pose = sbu.Pose(pose_estimator.estimate(im));
	// segmentation = sbu.Segmentation(im);
	// photo = sbu.Photo(im, 'pose', pose, 'segmentation', segmentation);

	im := ReadImage("file_name.png")
	poseEstimator := DefaultPoseEstimator()
	pose := poseEstimator.Estimate(im)
	VisualizePoseEstimation(im, pose)
}

func ReadImage(filepath string) image.Image {
	return nil
}

func VisualizePoseEstimation(im image.Image, pose Pose) {
}
