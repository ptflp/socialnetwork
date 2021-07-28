package components

import (
	"gitlab.com/InfoBlogFriends/server/cache"
	"gitlab.com/InfoBlogFriends/server/config"
	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/providers"
	"gitlab.com/InfoBlogFriends/server/respond"
	"gitlab.com/InfoBlogFriends/server/session"
	"go.uber.org/zap"
)

type Componenter interface {
	Logger() *zap.Logger
	Responder() respond.Responder
	LogLevel() zap.AtomicLevel
	JWTKeys() *session.JWTKeys
	Email() email.Mailer
	Config() *config.Config
	Cache() cache.Cache
	SMS() providers.SMS
	Decoder() *decoder.Decoder
	Facebook() providers.Socials
	Google() providers.Socials
}

type Components struct {
	logger    *zap.Logger
	responder respond.Responder
	logLevel  zap.AtomicLevel
	jwtKeys   *session.JWTKeys
	email     *email.Client
	config    *config.Config
	cache     cache.Cache
	sms       providers.SMS
	decoder   *decoder.Decoder
	facebook  providers.Socials
	google    providers.Socials
}

func (c *Components) Logger() *zap.Logger {
	return c.logger
}

func (c *Components) Responder() respond.Responder {
	return c.responder
}

func (c *Components) LogLevel() zap.AtomicLevel {
	return c.logLevel
}

func (c *Components) JWTKeys() *session.JWTKeys {
	return c.jwtKeys
}

func (c *Components) Email() email.Mailer {
	return c.email
}

func (c *Components) Config() *config.Config {
	return c.config
}

func (c *Components) Cache() cache.Cache {
	return c.cache
}

func (c *Components) SMS() providers.SMS {
	return c.sms
}

func (c *Components) Decoder() *decoder.Decoder {
	return c.decoder
}

func (c *Components) Facebook() providers.Socials {
	return c.facebook
}

func (c *Components) Google() providers.Socials {
	return c.google
}

func NewComponents(logger *zap.Logger) *Components {
	responder, err := respond.NewResponder(logger)
	if err != nil {
		logger.Fatal("responder initialization error", zap.Error(err))
	}

	// config initialization
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal("config initialization error", zap.Error(err))
	}

	c, err := cache.NewRedisCache(conf.Redis)
	if err != nil {
		logger.Fatal("redis initialization error", zap.Error(err))
	}

	jwt, err := session.NewJWTKeys(logger, c)
	if err != nil {
		logger.Fatal("jwt initialization error", zap.Error(err))
	}

	facebook := providers.NewFacebookAuth(&conf.Oauth2.Facebook)
	google := providers.NewGoogleAuth(&conf.Oauth2.Google)

	mailClient := email.NewClient(&conf.Email, logger)
	smsc := providers.NewSMSC(&conf.SMSC)

	return &Components{
		logger:    logger,
		responder: responder,
		logLevel:  zap.AtomicLevel{},
		jwtKeys:   jwt,
		email:     mailClient,
		config:    conf,
		cache:     c,
		sms:       smsc,
		decoder:   decoder.NewDecoder(),
		facebook:  facebook,
		google:    google,
	}
}
