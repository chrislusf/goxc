package tasks

/*
   Copyright 2013 Am Laher

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

import (
	//Tip for Forkers: please 'clone' from my url and then 'pull' from your url. That way you wont need to change the import path.
	//see https://groups.google.com/forum/?fromgroups=#!starred/golang-nuts/CY7o2aVNGZY
	"github.com/laher/goxc/archive"
	"github.com/laher/goxc/config"
	"github.com/laher/goxc/core"
	"github.com/laher/goxc/platforms"
	"log"
	"path/filepath"
)

// NOTE: in future this task may produce preferred types of archive for each OS (e.g. .tar.gz for Linux)
// TaskSettings should dictate this behaviour.

//runs automatically
func init() {
	Register(Task{
		"archive",
		"Create a compressed archive. Currently 'zip' format is used for all platforms except Linux (tar.gz)",
		runTaskArchive,
		map[string]interface{}{"os": map[string]interface{}{platforms.LINUX: "TarGz"}}})
}

func runTaskArchive(tp TaskParams) error {
	for _, dest := range tp.DestPlatforms {
		err := archivePlat(dest.Os, dest.Arch, tp.AppName, tp.WorkingDirectory, tp.OutDestRoot, tp.Settings)
		if err != nil {
			//TODO - 'force' option
			//return err
		}
	}
	//TODO return error?
	return nil
}

func archivePlat(goos, arch, appName, workingDirectory, outDestRoot string, settings config.Settings) error {
	resources := core.ParseIncludeResources(workingDirectory, settings.Resources.Include, settings.Resources.Exclude, settings.IsVerbose())

	// Create ZIP archive.
	relativeBin := core.GetRelativeBin(goos, arch, appName, false, settings.GetFullVersionName())

	var archiver archive.Archiver
	var ending string
	osOptions := settings.GetTaskSettingMap(TASK_ARCHIVE, "os")

	if osOption, keyExists := osOptions[goos]; keyExists {
		if osOption == "TarGz" {
			//if goos == core.LINUX {
			ending = "tar.gz"
			archiver = archive.TarGz
		} else {
			ending = "zip"
			archiver = archive.Zip
		}
	} else {
		ending = "zip"
		archiver = archive.Zip
	}

	archivePath, err := archive.ArchiveBinaryAndResources(
		filepath.Join(outDestRoot, settings.GetFullVersionName(), goos+"_"+arch),
		filepath.Join(outDestRoot, relativeBin), appName, resources, settings, archiver, ending)
	if err != nil {
		log.Printf("ZIP error: %s", err)
		return err
	} else {
		log.Printf("Artifact %s archived to %s", relativeBin, archivePath)
	}
	return nil
}
