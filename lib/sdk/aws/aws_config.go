package sdkaws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"os"
)

/**
 * aws_config
 * <p>
 * This file contains core data structures and logic used throughout the application.
 *
 * <p><strong>Copyright © 2025 – All rights reserved.</strong></p>
 *
 * <p>This source code is distributed under a collaborative license.</p>
 *
 * <p>
 * Contributions, suggestions, and improvements are welcome!
 * You are free to fork, modify, and submit pull requests under the terms of the repository's license.
 * Please ensure proper attribution to the original author(s) and preserve this notice in derivative works.
 * </p>
 *
 * @author Christian Bacilio De La Cruz
 * @email dbacilio88@outlook.es
 * @since 7/10/2025
 */

type AwsConfigParams struct {
	AccessKey string
	SecretKey string
	Region    string
	Profile   string
}

type AwsParams struct {
	Context context.Context
	Config  aws.Config
	Bucket  string
}

func ConfigAws(input AwsConfigParams) (aws.Config, error) {

	region := input.Region
	if region == "" {
		region = os.Getenv("REGION")
	}

	fmt.Println(region)
	profile := input.Profile
	if profile == "" {
		profile = os.Getenv("AWS_PROFILE")
	}
	fmt.Println(profile)
	var cfg aws.Config
	var err error

	var opts []func(*config.LoadOptions) error

	if input.AccessKey != "" && input.SecretKey != "" {
		cred := aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(input.AccessKey, input.SecretKey, ""))
		opts = append(opts, config.WithCredentialsProvider(cred))
	} else if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err = config.LoadDefaultConfig(context.TODO(), opts...)

	if err != nil {
		message = fmt.Sprintf("error loading config: %v", err)
		return aws.Config{}, errors.New(message)
	}

	fmt.Println("config loaded successfully")
	return cfg, nil
}
