# Speedcube events

It's a data scrapper from [worldcubeassociation.org](https://worldcubeassociation.org) and [speedcubing.pl/kalendarz-imprez](https://www.speedcubing.pl/kalendarz-imprez). It scraps Polish Speedcube events for 2024. Data is stored in S3 and presented on [https://kalendarz.krzysztofromanowski.pl/](https://kalendarz.krzysztofromanowski.pl/)

## Terraform RBAC

To be able to run jobs in Terraform, some things has to be added:

```bash
kubectl create clusterrole job_runner --verb=get,list,watch,create,update,patch,delete --resource=jobs,jobs/status

kubectl create clusterrolebinding job_runner_binding --clusterrole=job_runner --serviceaccount=default:default
```
