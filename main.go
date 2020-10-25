package main

//--------------------------------------------------------------------------------------------------
// PyEnv
//
// Pyenv is a tool used to install and remove Python versions, and to create and destroy Python
// virtual environments.
//
// Pyenv is MacOS-only.
//
// The goal of this program is to produce a tarball with the pyenv executable in it.  We
// do that by creating a hidden installation of homebrew, and use homebrew to install
// pyenv, also in a way that is hidden to the rest of the laptop software.  Then we create a
// tarball of the hidden homebrew directory.
//
//--------------------------------------------------------------------------------------------------

import (
  "fmt"
  "os"
  "os/exec"
  "os/user"
  _filepath "path/filepath"
)

type Builder struct {
}

func (this Builder) build() {
  repoDir, _ :=  _filepath.Abs(_filepath.Dir(os.Args[0]))
  pyenvWinDir := _filepath.Join(repoDir, "pyenv-win")
  currentUser, _ := user.Current();
  homeDir := currentUser.HomeDir    
  tarballDir := _filepath.Join(homeDir, "toolbox-tarballs")
  tarballFilepath := _filepath.Join(tarballDir, "toolbox-pyenv-win.tgz")

  // Find the git executable
  gitExecutable, err := exec.LookPath("git")
  if err != nil {
    fmt.Printf("error searching for 'git' executable. %s\n", err.Error())
    os.Exit(1)
  }

  // Find the tar executable
  tarExecutable, err := exec.LookPath("tar")
  if err != nil {
    fmt.Printf("error searching for 'tar' executable. %s\n", err.Error())
    os.Exit(1)
  }

  // Delete the pyenv-win directory if it exists
  _, err = os.Stat(pyenvWinDir)
  if !os.IsNotExist(err) {
    err = os.RemoveAll(pyenvWinDir)
    if err != nil {
      fmt.Printf("error deleting pyenv-win directory")
      os.Exit(1)
    }
  }

  // git clone
  command := exec.Command(gitExecutable, "clone", "https://github.com/pyenv-win/pyenv-win.git")
  fmt.Printf("--- %s\n", command.String());
  output, err := command.CombinedOutput();
  if err != nil {
    fmt.Printf("error running command '%s' '%s'. %s\n", command.String(), output, err.Error())
    os.Exit(1)
  }

  // Create the tarball directory if it does not exist
  _, err = os.Stat(tarballDir)
  if os.IsNotExist(err) {
    os.MkdirAll(tarballDir, 0700)
  }

  // tar the pyenv-win directory
  command = exec.Command(tarExecutable, "-czf", tarballFilepath, "pyenv-win")
  fmt.Printf("--- %s\n", command.String());
  output, err = command.CombinedOutput();
  if err != nil {
    fmt.Printf("error running command '%s' '%s'. %s\n", command.String(), output, err.Error())
    os.Exit(1)
  }
}

func main() {
  var builder Builder
  builder.build()
}
