$ sudo systemctl restart postgresql
$ sudo -i -u postgres

cloudflared tunnel --url tcp://localhost:15432 --origincert ~/.cloudflared/cert.pem --hostname pga1.schulenburg-office.com

cloudflared tunnel --config /dev/null --origincert ~/.cloudflared/cert.pem \
--url tcp://localhost:15432 --hostname pga1.schulenburg-office.com


## Cloudflare und pgAdmin
#### auf dem server
- $ cloudflared tunnel route dns schulenburg-office-tunnel pga1.schulenburg-office.com
  - $ nano /etc/cloudflared/config.yml    
    `- hostname: pga1.schulenburg-office.com`             
    `  service: tcp://localhost:5432`    
    
#### auf dem client
- $ cloudflared access tcp --hostname pga1.schulenburg-office.com --listener localhost:15432
- Im pgAdmin connection string localhost:15432 als host und port verwenden


#### DB dump von server auf client
1. Auf ssh server via
   $ ssh schulenburg-office    
   $ sudo -u postgres pg_dump trader > trader_dump.sql
   $ exit
2. Auf client
   $ scp schulenburg-office:~/trader_dump.sql .
3. $ sudo cp trader_dump.sql /tmp/
4. $ sudo chmod 644 /tmp/trader_dump.sql
5. pgAdmin   
   drop table chart;    
   drop table tradingPlan;    
   drop table "user";    
   drop table transaction;    
   drop table movingAverage;    
6. sudo -u postgres psql -d trader -f /tmp/trader_dump.sql



SELECT tp.id, tp.asset, tp.startPips, tp.endPips, tp.startDollar, tp.endDollar, tp.startEuro, tp.endEuro, tp.startAt,
 tp.endAt,
 c1.close as "lastPips", c2.close as "eurUsd", c1.chartAt as "lastPipsAt", tp.deviseStart, tp.deviseEnd, c3.deviseLast, c3.deviseLastAt, tp.deviseCur, tp.deviseCurUsd, tp.created, tp.modified, c1.sourceService
 FROM tradingPlan tp
 LEFT JOIN LATERAL (
  SELECT c10.close, c10.chartAt, c10.sourceService
     FROM chart c10
     WHERE c10.asset = tp.asset
     ORDER BY c10.chartAt DESC
     LIMIT 1
 ) c1 ON true
 LEFT JOIN LATERAL (
     SELECT c11.close
     FROM chart c11
     WHERE c11.asset = 'EURUSD'
     ORDER BY c11.chartAt DESC
     LIMIT 1
 ) c2 ON true
 LEFT JOIN LATERAL (
     SELECT c13.close as deviseLast, c13.chartAt as deviseLastAt
     FROM chart c13
     WHERE c13.asset = tp.deviseCurUsd
     ORDER BY c13.chartAt DESC
     LIMIT 1
 ) c3 ON true
 WHERE endAt IS NOT NULL
 order by tp.id DESC;

