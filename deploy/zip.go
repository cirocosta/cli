package deploy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/configstore"
	"github.com/launchpad-project/cli/verbose"
	"github.com/sabhiram/go-git-ignore"
)

type JarContainer struct {
	ContainerPath string
	Writer        *zip.Writer
	strings       []string
	IgnoreRules   *ignore.GitIgnore
}

var GlobalIgnoreList = []string{
	".DS_Store",
	".directory",
	".Trashes",
	".project",
	".settings",
	".idea",
}

func WriteJar(dest, dirPath string, ignoreList []string) error {
	ignoreList = append(ignoreList, GlobalIgnoreList...)
	ignoreRules, err := ignore.CompileIgnoreLines(ignoreList...)

	if err != nil {
		return err
	}

	z, err := os.Create(dest)

	if err != nil {
		return err
	}

	jar := &JarContainer{
		ContainerPath: dirPath,
		Writer:        zip.NewWriter(z),
		IgnoreRules:   ignoreRules,
	}

	err = filepath.Walk(dirPath, jar.walkFunc)

	if err == nil {
		err = jar.Writer.Close()
	}

	if err == nil {
		err = z.Close()
	}

	return err
}

func ZipContainer(container, out string) error {
	var dir = filepath.Join(config.Context.ProjectRoot, container)

	var cs = &configstore.Store{
		Name: container,
		Path: filepath.Join(dir, "container.json"),
	}

	err := cs.Load()

	if err == nil {
		err = WriteJar(out, dir, getIgnoreList(cs))
	}

	return err
}

func ZipContainers(list []string) {
	for _, container := range list {
		fmt.Println("Building", container, "JAR bundle")

		err := ZipContainer(container, container+".jar")

		if err != nil {
			println(err)
			os.Exit(1)
		}
	}
}

func (j *JarContainer) getConditionalSkipDir(fi os.FileInfo) error {
	if fi.IsDir() {
		// @todo optimize avoiding going inside folders
		// that obviously are totally ignored
		// maybe: "when no string with ! or * exists, use return filepath.SkipDir instead"
		// return filepath.SkipDir
	}

	return nil
}

func (j *JarContainer) walkFunc(path string, fi os.FileInfo, ierr error) error {
	if ierr != nil {
		verbose.Debug("Error walking to " + path)
		panic(ierr)
	}

	relative, err := filepath.Rel(j.ContainerPath, path)

	if err != nil {
		panic(err)
	}

	if relative == "." {
		return nil
	}

	header, err := zip.FileInfoHeader(fi)

	if err != nil {
		panic(err)
	}

	if j.IgnoreRules.MatchesPath(relative) {
		return j.getConditionalSkipDir(fi)
	}

	header.Name = relative

	if fi.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	writer, err := j.Writer.CreateHeader(header)

	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	file, err := os.Open(path)

	if err == nil {
		if verbose.Enabled {
			stat, _ := file.Stat()
			verbose.Debug(fmt.Sprintf("%v (%v bytes)", relative, stat.Size()))
		}
		_, err = io.Copy(writer, file)
	}

	if err == nil {
		err = file.Close()
	}

	return err
}
