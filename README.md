# create_dify_workspace

rodar esse comando para descobrir o ip do db e alterar no .env
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' automations_dify-db-1

após rodar o script entrar no container da api, entrar no /storage/privkey

e criar a nova privkey com o id do tenants