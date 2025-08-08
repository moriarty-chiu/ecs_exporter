### listBuckets
```python
from obs import ObsClient
import os
import traceback

# 推荐通过环境变量获取AKSK，这里也可以使用其他外部引入方式传入,如果使用硬编码可能会存在泄露风险
# 您可以登录访问管理控制台获取访问密钥AK/SK，获取方式请参见https://support.huaweicloud.com/usermanual-ca/ca_01_0003.html
ak = os.getenv("AccessKeyID")
sk = os.getenv("SecretAccessKey")
# 【可选】如果使用临时AKSK和SecurityToken访问OBS，则同样推荐通过环境变量获取
# security_token = os.getenv("SecurityToken")
# server填写Bucket对应的Endpoint, 这里以华北-北京四为例，其他地区请按实际情况填写
server = "https://obs.cn-north-4.myhuaweicloud.com"
# 创建obsClient实例
# 如果使用临时AKSK和SecurityToken访问OBS，需要在创建实例时通过security_token参数指定securityToken值
obsClient = ObsClient(access_key_id=ak, secret_access_key=sk, server=server)
try:
    # 列举桶，并设置isQueryLocation参数为True，同时查询桶区域
    resp = obsClient.listBuckets(True)
    # 返回码为2xx时，接口调用成功，否则接口调用失败
    if resp.status < 300:
        print('List Buckets Succeeded')
        print('requestId:', resp.requestId)
        print('name:', resp.body.owner.owner_id)
        print('create_date:', resp.body.owner.owner_name)
        index = 1
        for bucket in resp.body.buckets:
            print('bucket [' + str(index) + ']')
            print('name:', bucket.name)
            print('create_date:', bucket.create_date)
            print('location:', bucket.location)
            index += 1
    else:
        print('List Buckets Failed')
        print('requestId:', resp.requestId)
        print('errorCode:', resp.errorCode)
        print('errorMessage:', resp.errorMessage)
except:
    print('List Buckets Failed')
    print(traceback.format_exc())
```
### getBucketStorageInfo
```python
from obs import ObsClient
import os
import traceback

# 推荐通过环境变量获取AKSK，这里也可以使用其他外部引入方式传入，如果使用硬编码可能会存在泄露风险
# 您可以登录访问管理控制台获取访问密钥AK/SK，获取方式请参见https://support.huaweicloud.com/usermanual-ca/ca_01_0003.html
ak = os.getenv("AccessKeyID")
sk = os.getenv("SecretAccessKey")
# 【可选】如果使用临时AKSK和SecurityToken访问OBS，则同样推荐通过环境变量获取
# security_token = os.getenv("SecurityToken")
# server填写Bucket对应的Endpoint, 这里以华北-北京四为例，其他地区请按实际情况填写
server = "https://obs.cn-north-4.myhuaweicloud.com"
# 创建obsClient实例
# 如果使用临时AKSK和SecurityToken访问OBS，需要在创建实例时通过security_token参数指定securityToken值
obsClient = ObsClient(access_key_id=ak, secret_access_key=sk, server=server)
try:
    bucketName="examplebucket"
    #获取桶的存量信息
    resp = obsClient.getBucketStorageInfo(bucketName)
    # 返回码为2xx时，接口调用成功，否则接口调用失败
    if resp.status < 300:
        print('Get Bucket StorageInfo Succeeded')
        print('requestId:', resp.requestId)
        print('size:', resp.body.size)
        print('objectNumber:', resp.body.objectNumber)
    else:
        print('Get Bucket StorageInfo Failed')
        print('requestId:', resp.requestId)
        print('errorCode:', resp.errorCode)
        print('errorMessage:', resp.errorMessage)
except:
    print('Get Bucket StorageInfo Failed')
    print(traceback.format_exc())
```
### getBucketQuota
```python
from obs import ObsClient
import os
import traceback

# 推荐通过环境变量获取AKSK，这里也可以使用其他外部引入方式传入，如果使用硬编码可能会存在泄露风险
# 您可以登录访问管理控制台获取访问密钥AK/SK，获取方式请参见https://support.huaweicloud.com/usermanual-ca/ca_01_0003.html
ak = os.getenv("AccessKeyID")
sk = os.getenv("SecretAccessKey")
# 【可选】如果使用临时AKSK和SecurityToken访问OBS，则同样推荐通过环境变量获取
# security_token = os.getenv("SecurityToken")
# server填写Bucket对应的Endpoint, 这里以华北-北京四为例，其他地区请按实际情况填写
server = "https://obs.cn-north-4.myhuaweicloud.com"
# 创建obsClient实例
# 如果使用临时AKSK和SecurityToken访问OBS，需要在创建实例时通过security_token参数指定securityToken值
obsClient = ObsClient(access_key_id=ak, secret_access_key=sk, server=server)
try:
    bucketName="examplebucket"
    #获取桶配额
    resp = obsClient.getBucketQuota(bucketName)
    # 返回码为2xx时，接口调用成功，否则接口调用失败
    if resp.status < 300:
        print('Get Bucket Quota Succeeded')
        print('requestId:', resp.requestId)
        print('quota:', resp.body.quota)
    else:
        print('Get Bucket Quota Failed')
        print('requestId:', resp.requestId)
        print('errorCode:', resp.errorCode)
        print('errorMessage:', resp.errorMessage)
except:
    print('Get Bucket Quota Failed')
    print(traceback.format_exc())
```