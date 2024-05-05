swag init -o api/rest/docs
if [ $? -eq 0 ]
then
  echo "bcapi swagger success"
else
  echo "bcapi swagger failed" >&2
  exit 1
fi


GOOS=linux GOARCH=amd64 go build  -o bcwork main.go;
if [ $? -eq 0 ] 
then 
  echo "bcwork build success"
else 
  echo "bcwork build failed" >&2
  exit 1
fi

mv bcwork deploy/ansible/files/
ansible-playbook -u root -i /etc/ansible/brightcom_hosts deploy/ansible/bcwork_update.yml

