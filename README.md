# eventus
A wrapper library in GO for cloud messaging


## Goal
To bring different messaging service providers under single shell and provide an interface to use them as per requirement in event driven architectures.


## Pre-requsite
- GO version **1.13 and above**
- While using GCP provider **GOOGLE_APPLICATION_CREDENTIALS** should be set as an environment variable with value as path to project config json which can obtained from Google Cloud Account
- For AWS, *.aws* folder on host machine should contain **credentials** file associated with user's AWS account


## Execution
  ```bash
  go run example/main.go 
  ``` 
  Run above command to execute example file. Make sure to replace topics, IDs with valid values.
  

## License
[MIT](https://choosealicense.com/licenses/mit/)
