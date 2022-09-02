package oauth

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

// Config oauth2.Config overridden with device endpoint for supporting Device Authorization Grant flow [RFC8268]
type Config struct {
	*oauth2.Config
	DeviceEndpoint string
}

// BuildConfigFromMetadataService builds OAuth2 config from information retrieved through the anonymous auth metadata service.
func BuildConfigFromMetadataService(ctx context.Context, authMetadataClient service.AuthMetadataServiceClient) (clientConf *Config, err error) {
	var clientResp *service.PublicClientAuthConfigResponse
	if clientResp, err = authMetadataClient.GetPublicClientConfig(ctx, &service.PublicClientAuthConfigRequest{}); err != nil {
		return nil, err
	}

	var oauthMetaResp *service.OAuth2MetadataResponse
	if oauthMetaResp, err = authMetadataClient.GetOAuth2Metadata(ctx, &service.OAuth2MetadataRequest{}); err != nil {
		return nil, err
	}

	clientConf = &Config{
		Config: &oauth2.Config{
			ClientID:    clientResp.ClientId,
			RedirectURL: clientResp.RedirectUri,
			Scopes:      clientResp.Scopes,
			Endpoint: oauth2.Endpoint{
				TokenURL: oauthMetaResp.TokenEndpoint,
				AuthURL:  oauthMetaResp.AuthorizationEndpoint,
			},
		},
		DeviceEndpoint: oauthMetaResp.DeviceAuthorizationEndpoint,
	}

	return clientConf, nil
}
