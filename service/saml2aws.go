package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Provider struct {
	filePath    string
	profilePath string
}

type MFA struct {
	Token string
}

type LoginOption struct {
	Console       bool
	LinkOnly      bool
	Firefox       bool
	UpdateProfile bool
}

func GetSAML2AWS() (*Provider, error) {
	// Check saml2aws
	out, err := exec.Command("saml2aws").Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		return nil, err
	}

	home := os.Getenv("HOME")
	filePath := fmt.Sprintf("%s/.saml2aws", home)
	profilePath := fmt.Sprintf("%s/.awsctx", home)

	return &Provider{filePath, profilePath}, nil
}

func (pr *Provider) Login(profile string, opts *LoginOption, mfa *MFA) error {
	opt := "login"
	if opts.Console {
		opt = "console"
	}

	baseCommand := []string{opt, "-a", profile, "--skip-prompt", "--quiet"}
	if mfa != nil {
		baseCommand = append(baseCommand, "--mfa-token", mfa.Token)
	}

	// Link Only
	if opts.LinkOnly || opts.Firefox {
		baseCommand = append(baseCommand, "--link")
	}

	if !opts.LinkOnly || (opts.LinkOnly && opts.Firefox) {
		defer fmt.Println(profile)
	}

	// Firefox
	fxBin := ""
	if opts.Firefox {
		if f := os.Getenv("FIREFOX_BIN"); len(f) != 0 {
			if err := exec.Command(f, "--help").Run(); err == nil {
				fxBin = f
			}
		} else if err := exec.Command("firefox", "--help").Run(); err == nil {
			fxBin = "firefox"
		} else if err := exec.Command("firefox.exe", "--help").Run(); err == nil {
			fxBin = "firefox.exe"
		} else if err := exec.Command("/Applications/Firefox.app/Contents/MacOS/firefox", "--help").Run(); err == nil {
			fxBin = "/Applications/Firefox.app/Contents/MacOS/firefox"
		}

		if fxBin == "" {
			fmt.Println("Firefox not found. You can set $FIREFOX_BIN manually.")
			return errors.New("could not find firefox")
		}
	}

	// Call saml2aws
	cmd := exec.Command("saml2aws", baseCommand...)
	res, err := cmd.Output()
	if opts.LinkOnly {
		fmt.Println(string(res))
	}
	if err != nil {
		return err
	}

	// Open firefox
	if opts.Firefox {
		color := "turquoise"
		if strings.Contains(profile, "develop") || strings.Contains(profile, "dev") {
			color = "green"
		} else if strings.Contains(profile, "staging") || strings.Contains(profile, "stg") {
			color = "yellow"
		} else if strings.Contains(profile, "production") || strings.Contains(profile, "prod") {
			color = "pink"
		}

		err := pr.execute(fxBin, fmt.Sprintf(
			"ext+container:name=%s&color=%s&icon=chill&url=%s",
			profile,
			color,
			strings.ReplaceAll(string(res), "&", "%26"),
		)).Run()
		if err != nil {
			return err
		}
	}

	// Update current profile
	if opts.UpdateProfile {
		f, err := os.Create(pr.profilePath)
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(profile))
		if err != nil {
			return err
		}

		// Update env
		os.Setenv("AWS_PROFILE", profile)
	}
	return nil
}

func (pr *Provider) ListRoles(profile string, mfa *MFA) error {
	baseCommand := []string{"list-roles", "-a", profile, "--skip-prompt"}
	if mfa != nil {
		baseCommand = append(baseCommand, "--mfa-token", mfa.Token)
	}

	cmd := pr.execute("saml2aws", baseCommand...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (pr *Provider) GetProfiles() ([]string, error) {
	file, err := os.Open(pr.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	profiles := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profiles = append(profiles, strings.TrimPrefix(strings.TrimSuffix(line, "]"), "["))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (pr *Provider) execute(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
