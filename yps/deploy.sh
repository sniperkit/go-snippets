#!/bin/bash
appcfg.py --oauth2 update default/app.yaml backend/backend.yaml processing/processing.yaml
appcfg.py --oauth2 update_dispatch default
appcfg.py --oauth2 update_queues default
appcfg.py --oauth2 update_cron default
