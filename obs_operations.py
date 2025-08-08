#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Huawei Cloud OBS SDK Python Example
Retrieves information for ALL buckets instead of just the first one
"""

import os
import traceback
from obs import ObsClient


def initialize_obs_client():
    """
    Initialize OBS client with credentials from environment variables
    """
    # Recommended to get AKSK via environment variables to avoid leakage risk
    # You can get access keys by logging into the management console
    # See: https://support.huaweicloud.com/usermanual-ca/ca_01_0003.html
    ak = os.getenv("AccessKeyID")
    sk = os.getenv("SecretAccessKey")
    
    # Server should be set to the endpoint corresponding to your bucket
    # This example uses North China - Beijing 4
    server = "https://obs.cn-north-4.myhuaweicloud.com"
    
    # If using temporary AKSK and SecurityToken, also recommended to get via environment variables
    # security_token = os.getenv("SecurityToken")
    
    try:
        # Create ObsClient instance
        # If using temporary AKSK and SecurityToken, specify securityToken via security_token parameter
        obs_client = ObsClient(access_key_id=ak, secret_access_key=sk, server=server)
        print("OBS client initialized successfully")
        return obs_client
    except Exception as e:
        print(f"Failed to initialize OBS client: {e}")
        return None


def list_buckets(obs_client):
    """
    List all buckets and their locations
    """
    try:
        # List buckets and set isQueryLocation parameter to True to also query bucket regions
        resp = obs_client.listBuckets(isQueryLocation=True)
        
        # When return code is 2xx, the API call is successful, otherwise it fails
        if resp.status < 300:
            print('List Buckets Succeeded')
            print('requestId:', resp.requestId)
            print('owner_id:', resp.body.owner.owner_id)
            print('owner_name:', resp.body.owner.owner_name)
            
            buckets = resp.body.buckets
            if buckets:
                print(f'Total buckets found: {len(buckets)}')
                index = 1
                for bucket in buckets:
                    print(f'bucket [{index}]')
                    print('name:', bucket.name)
                    print('create_date:', bucket.create_date)
                    print('location:', bucket.location)
                    index += 1
                return buckets
            else:
                print("No buckets found")
                return []
        else:
            print('List Buckets Failed')
            print('requestId:', resp.requestId)
            print('errorCode:', resp.errorCode)
            print('errorMessage:', resp.errorMessage)
            return None
    except Exception as e:
        print('List Buckets Failed')
        print(traceback.format_exc())
        return None


def get_bucket_storage_info(obs_client, bucket_name):
    """
    Get storage information of a bucket
    """
    try:
        # Get bucket storage information
        resp = obs_client.getBucketStorageInfo(bucket_name)
        
        # When return code is 2xx, the API call is successful, otherwise it fails
        if resp.status < 300:
            print(f'Get StorageInfo for Bucket "{bucket_name}" Succeeded')
            print('  requestId:', resp.requestId)
            print('  size:', resp.body.size, 'bytes')
            print('  objectNumber:', resp.body.objectNumber)
            return resp.body
        else:
            print(f'Get StorageInfo for Bucket "{bucket_name}" Failed')
            print('  requestId:', resp.requestId)
            print('  errorCode:', resp.errorCode)
            print('  errorMessage:', resp.errorMessage)
            return None
    except Exception as e:
        print(f'Get StorageInfo for Bucket "{bucket_name}" Failed')
        print(traceback.format_exc())
        return None


def get_bucket_quota(obs_client, bucket_name):
    """
    Get quota information of a bucket
    """
    try:
        # Get bucket quota
        resp = obs_client.getBucketQuota(bucket_name)
        
        # When return code is 2xx, the API call is successful, otherwise it fails
        if resp.status < 300:
            print(f'Get Quota for Bucket "{bucket_name}" Succeeded')
            print('  requestId:', resp.requestId)
            print('  quota:', resp.body.quota, 'bytes')
            return resp.body
        else:
            print(f'Get Quota for Bucket "{bucket_name}" Failed')
            print('  requestId:', resp.requestId)
            print('  errorCode:', resp.errorCode)
            print('  errorMessage:', resp.errorMessage)
            return None
    except Exception as e:
        print(f'Get Quota for Bucket "{bucket_name}" Failed')
        print(traceback.format_exc())
        return None


def main():
    """
    Main function demonstrating OBS operations on ALL buckets
    """
    # Initialize OBS client
    obs_client = initialize_obs_client()
    if not obs_client:
        print("Failed to initialize OBS client. Exiting.")
        return
    
    try:
        # List all buckets
        print("=== Listing Buckets ===")
        buckets = list_buckets(obs_client)
        
        if buckets:
            # Process storage info and quota for ALL buckets
            for bucket in buckets:
                bucket_name = bucket.name
                print(f"\n=== Getting Storage Info for Bucket: {bucket_name} ===")
                get_bucket_storage_info(obs_client, bucket_name)
                
                print(f"\n=== Getting Quota for Bucket: {bucket_name} ===")
                get_bucket_quota(obs_client, bucket_name)
        else:
            # If no buckets exist, try example bucket
            example_bucket = "examplebucket"
            print(f"\n=== Getting Storage Info for Example Bucket: {example_bucket} ===")
            get_bucket_storage_info(obs_client, example_bucket)
            
            print(f"\n=== Getting Quota for Example Bucket: {example_bucket} ===")
            get_bucket_quota(obs_client, example_bucket)
    
    except Exception as e:
        print(f"An error occurred during execution: {e}")
        print(traceback.format_exc())
    
    finally:
        # Close the OBS client connection
        if obs_client:
            obs_client.close()
            print("\nOBS client connection closed")


if __name__ == "__main__":
    main()