# This script use to wakeup terminated SPOT VMs
#
# Run the script by cron every 1 minute
# Run command crontab -e then add the line below
# * * * * * sh wakeup-vm.sh

# put list of instances here
for instance in "iplay-dev-api01"
do
  is_terminated=$(gcloud compute instances list | grep $instance | grep TERMINATED)

  if [ -z $is_terminated ]
  then
    echo "$instance did not terminate"
    exit
  fi

  echo "start $instance"
  gcloud compute instances start $instance --zone=us-central1-a
done
