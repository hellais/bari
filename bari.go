package main

import (
  "os"
  "log"
  "fmt"
  "strings"

  "github.com/codegangsta/cli"
)

func show_install_instructions(name string) {
  pkgs := loadPackages(name);

  for _, pkg := range pkgs {
    if (pkg.os != "") {
      output := "# " + pkg.os;
      if (pkg.distro != "") {
        output = output + " " + pkg.distro;
      }
      if (pkg.release != "") {
        output = output + " release: " + pkg.release;
      }
      fmt.Println(output);
    }
    
    if (pkg.pkg_manager != "") {
      output := "# using " + pkg.pkg_manager;
      if (pkg.pkg_manager_version != "") {
        output = output + " version: " + pkg.pkg_manager_version;
      }
      fmt.Println(output);
    }
    
    fmt.Println("# run:");
    fmt.Println(strings.Join(pkg.install_command(), " "));
    fmt.Println();
  }
}

func install_package(name string) {
  var platform Platform;
  platform.detect();

  var supportedPackages []Package;

  pkgs := loadPackages(name);

  for _, pkg := range pkgs {
    if platform.supports(pkg) {
      supportedPackages = append(supportedPackages, pkg);
    }
  }
  
  if len(supportedPackages) == 0 {
    log.Fatal("Could not find a supported install method for your OS");
  }
  
  fmt.Println("# Install options");
  for _, pkg := range supportedPackages {
    pkg_manager, ok := pkg.json["pkg_manager"];
    if (ok) {
      fmt.Println("##", pkg_manager);
    }
    pkg_manager_version, ok := pkg.json["pkg_manager_version"];
    if (ok) {
      fmt.Println("Package manager version:", pkg_manager_version);
    }

    OS, ok := pkg.json["os"];
    if (ok) {
      fmt.Println("Operating system:", OS);
    }

    release, ok := pkg.json["release"];
    if (ok) {
      fmt.Println("Release:", release);
    }

    pkg_name, ok := pkg.json["pkg"];
    if (ok) {
      fmt.Println("Package name:", pkg_name);
    }

    fmt.Println("Install command: ", pkg.install_command());
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "bari"
  app.Usage = "Install everything, everywhere."
  app.Commands = []cli.Command{
    {
      Name: "install",
      Aliases: []string{"i"},
      Usage: "Install the specified package on your system",
      Action: func(c *cli.Context) {
        pkg_name := c.Args().First();
        if pkg_name == "" {
          cli.ShowAppHelp(c);
          log.Fatal("You MUST specify a package name");
        }
        //log.Print("Installing ", pkg_name);
        install_package(pkg_name);
      },
    },

    {
      Name: "show",
      Usage: "Show the setup instructions for the specified package",
      Action: func(c *cli.Context) {
        pkg_name := c.Args().First();
        if pkg_name == "" {
          cli.ShowAppHelp(c);
          log.Fatal("You MUST specify a package name");
        }
        show_install_instructions(pkg_name);
      },
    },
  };
  app.Run(os.Args)
}
