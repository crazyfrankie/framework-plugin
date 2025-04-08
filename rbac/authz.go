package rbac

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	defaultPolicy = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act`
)

// Authz defines an authorizer
type Authz struct {
	*casbin.SyncedEnforcer
}

type option struct {
	aclModel           string        // Casbin's model string
	autoLoadPolicyTime time.Duration // Interval for auto-loading policies
}

type Option func(*option)

func WithACLModel(model string) Option {
	return func(o *option) {
		o.aclModel = model
	}
}

func WithLoadTime(interval time.Duration) Option {
	return func(o *option) {
		o.autoLoadPolicyTime = interval
	}
}

func defaultOption() *option {
	return &option{
		aclModel:           defaultPolicy,
		autoLoadPolicyTime: time.Second * 5,
	}
}

func NewAuthz(db *gorm.DB, opts ...Option) (*Authz, error) {
	opt := defaultOption()

	for _, o := range opts {
		o(opt)
	}

	// initialize the Gorm adapter and use it for the Casbin Authorizer
	adp, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// load Matching Policy Model from Configuration
	md, _ := model.NewModelFromString(opt.aclModel)

	enforcer, err := casbin.NewSyncedEnforcer(md, adp)
	if err != nil {
		return nil, err
	}

	// load permission data from db
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// start the autoloaded policy, using the configured time interval
	enforcer.StartAutoLoadPolicy(opt.autoLoadPolicyTime)

	return &Authz{SyncedEnforcer: enforcer}, nil
}

// Authorize for authorization
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	// Call the Enforce method for authorization checking
	return a.Enforce(sub, obj, act)
}
