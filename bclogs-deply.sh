GOOS=linux GOARCH=amd64 go build  -o bcwork main.go;
if [ $? -eq 0 ] 
then 
  echo "bclogs build success"
else 
  echo "bclogs build failed" >&2
  exit 1
fi

mv bcwork deploy/ansible/files/
#ansible-playbook -u admin deploy/ansible/bcwork_update.yml
ansible-playbook -u root deploy/ansible/bcwork_logs_update.yml

