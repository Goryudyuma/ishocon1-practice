echo "" > /var/log/mysql-slow-logs.log 
echo "" > sudo tee /var/log/nginx/access.log
sudo service nginx restart
sudo service mysqld restart
