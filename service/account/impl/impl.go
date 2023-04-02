package impl

import (
	account "actuator/service/account/api"
)

type AccountService struct {
	account.UnimplementedAccountServer
}
