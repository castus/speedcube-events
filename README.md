# Speedcube events

The root of _Speedcube events_ repos family.

The _Speedcube events_ project's purpose is to create a place when all Poland Speedcube Events can be displayed with details like: the main competition, all competitions, limit of competitors, currently registered competitors etc.

## How does it work?

1. The source of truth is [speedcubing.pl/kalendarz-imprez](https://www.speedcubing.pl/kalendarz-imprez) page, which is scraped and parsed.
1. The scraped data is then fill up with all details using 3rd party sites. Depends of what's the contest, it can be:
    * [worldcubeassociation.org](https://worldcubeassociation.org)
    * [cube4fun.pl](https://cube4fun.pl)
    * [rubiart.pl](https://rubiart.pl)
    * etc.
1. Data is presented on [https://kalendarz.krzysztofromanowski.pl/](https://kalendarz.krzysztofromanowski.pl/)

## Technical details

1. Data is stored in DynamoDB and exported as json to S3.
1. The page is hosted using Vercel
1. To scrap services like Cube4Fun or PPO, Selenium web driver is used.
   1. For Cube4Fun event, the WCA search mechanism is used first, and when there is a result, data is taken from there.
1. The core scrap mechanism is run in on-premise K8s using CronJobs.
1. Web scraping is spin up automatically as K8s Jobs.

## Repositories

* [castus/speedcube-events](https://github.com/castus/speedcube-events) - this repo, the core scrap functionality
* [castus/speedcube-events-charts](https://github.com/castus/speedcube-events-charts) - K8s charts to spin up the core service
* [castus/speedcube-events-app](https://github.com/castus/speedcube-events-app) - the webpage that displays the data
* [castus/speedcube-events-web-scraper](https://github.com/castus/speedcube-events-web-scraper) - image to be used for scrap webpages that uses JS to display

## Core scrap functionality

This repo produces a CLI binary. Use `--help` to see all available commands.

## Terraform RBAC

To be able to run jobs in Terraform, some things has to be added:

```bash
kubectl create clusterrole job_runner --verb=get,list,watch,create,update,patch,delete --resource=jobs,jobs/status
kubectl create clusterrolebinding job_runner_binding --clusterrole=job_runner --serviceaccount=default:default
```

## Support services

1. [Brevo](https://brevo.com) - SMTP email relay
2. [Amplitude](https://amplitude.com) - Analytics
