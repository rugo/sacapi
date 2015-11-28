#! /bin/sh
### BEGIN INIT INFO
# Provides:          saccalendar
# Required-Start:    $syslog $time $remote_fs
# Required-Stop:     $syslog $time $remote_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Smart Alarm Clock API
# Description:       Debian init script for the Smart Alarm Clock API
### END INIT INFO
#
# Author:	RuGo rg@ht11.org
#

PATH=/bin:/usr/bin:/sbin:/usr/sbin:/home/rg/go/bin
DAEMON=$(which saccalendar)
PIDFILE=/tmp/saccalendar.pid

test -x $DAEMON || exit 0

. /lib/lsb/init-functions

case "$1" in
  start)
	log_daemon_msg "Starting saccalendar" "saccalendar"
	start_daemon -p $PIDFILE $DAEMON &
	log_end_msg $?
    ;;
  stop)
	log_daemon_msg "Stopping saccalendar" "saccalendar"
	killproc -p $PIDFILE $DAEMON
	log_end_msg $?
    ;;
  force-reload|restart)
    $0 stop
    $0 start
    ;;
  status)
    status_of_proc -p $PIDFILE $DAEMON atd && exit 0 || exit $?
    ;;
  *)
    echo "Usage: /etc/init.d/atd {start|stop|restart|force-reload|status}"
    exit 1
    ;;
esac

exit 0
