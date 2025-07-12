package sdkaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/**
 * aws_s3
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

var message string

func ListFilesBucket(input AwsParams) ([]string, error) {

	var fileNames []string
	client := s3.NewFromConfig(input.Config)

	params := s3.ListObjectsV2Input{
		Bucket: aws.String(input.Bucket),
	}

	list, err := client.ListObjectsV2(input.Context, &params)
	if err != nil {
		message = fmt.Sprintf("error list to files in bucket s3: %v", err)
		return fileNames, fmt.Errorf(message)
	}

	for _, file := range list.Contents {
		fileNames = append(fileNames, *file.Key)
	}

	return fileNames, nil

}
