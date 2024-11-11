// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

type KeystoreModel struct {
	Password   string
	Base64Text string
}

const (
	FILENAME  = "79e00021-58b2-4652-b498-59134c0ed6e7.pkcs12"
	FILENAME2 = "79e00021-58b2-4652-b498-59134c0ed6e72.pkcs12"
)

func (m KeystoreModel) CreateKeystoreBase64() (string, error) {
	fileName := FILENAME

	cmd := exec.Command(
		"keytool",
		"-v",
		"-genkeypair",
		"-alias", "keystore",
		"-keypass", m.Password,
		"-keystore", fileName,
		"-storepass", m.Password,
		"-validity", "10000",
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-dname", "CN=Unknown, OU=Unknown, O=Unknown, L=Unknown, S=Unknown, C=Unknown",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error after executing command to produce keystore file. Is the keytool installed on the machine?\nError:%s\nOutput: %s", err, string(output))
	}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("error reading keystore file after producing keystore file\nError: %s", err)
	}
	err = os.Remove(fileName)
	if err != nil {
		return "", fmt.Errorf("error removing keystore file after reading keystore file\nError: %s", err)
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (oldModel KeystoreModel) UpdateKeystoreBase64(newPassword string) (string, error) {

	fileName := FILENAME
	fileName2 := FILENAME2

	decoded, err := base64.StdEncoding.DecodeString(oldModel.Base64Text)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fileName, decoded, 0644); err != nil {
		return "", err
	}

	changePassCmd := exec.Command(
		"keytool",
		"-importkeystore",
		"-srckeystore", fileName,
		"-srcstoretype", "PKCS12",
		"-srcstorepass", oldModel.Password,
		"-destkeystore", fileName2,
		"-deststoretype", "PKCS12",
		"-deststorepass", newPassword,
		"-destkeypass", newPassword,
	)

	if out, err := changePassCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("error changing password\nError: %s\nOutput: %s", err, string(out))
	}

	if err := os.Remove(fileName); err != nil {
		return "", fmt.Errorf("error removing old keystore file after changing password\nFile name: %s\nError: %s", fileName, err)
	}

	bytes, err := os.ReadFile(fileName2)
	if err != nil {

		return "", fmt.Errorf("error reading new keystore file after changing password\nFile name: %s\nError: %s", fileName2, err)
	}

	if err := os.Remove(fileName2); err != nil {
		return "", fmt.Errorf("error removing keystore file after reading keystore file. Error: %s", err)
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
