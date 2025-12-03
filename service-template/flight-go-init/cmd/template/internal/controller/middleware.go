package controller

import (
	"context"
	"net/http"
	"time"
	"{{MODULE_NAME}}/internal/shared/deps"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

func HawkEyeAPIInbound(deps deps.Deps) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := time.Now()
			tags := make(map[string]interface{})
			tags["protocol"] = "REST"

			err := next(c)
			if err != nil {
				go func(t time.Time) {
					err := deps.Metric.CustomMonitorLatency(c.Request().URL.RequestURI(), metrics.API_IN,
						metrics.Failed, http.StatusInternalServerError, tags, time.Since(t))
					if err != nil {
						deps.GetLogger(nil).Errorf("%s - Error publish to statsd (hawkeye). Got Err: %s", "[StatsD]",
							err.Error())
					}
				}(t)

				return err
			}

			go func(t time.Time) {
				err := deps.Metric.CustomMonitorLatency(c.Request().URL.RequestURI(), metrics.API_IN,
					metrics.Success, http.StatusOK, tags, time.Since(t))
				if err != nil {
					deps.GetLogger(nil).Errorf("%s - Error publish to statsd (hawkeye). Got Err: %s", "[StatsD]",
						err.Error())
				}
			}(t)
			return nil
		}
	}
}

func MandatoryRequestMiddleware(deps deps.Deps) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			storeID := c.Request().Header.Get("storeId")
			channelID := c.Request().Header.Get("channelId")
			requestID := c.Request().Header.Get("requestId")
			serviceID := c.Request().Header.Get("serviceId")
			resellerID := c.Request().Header.Get("resellerId")
			identity := c.Request().Header.Get("identity")
			username := c.Request().Header.Get("username")

			XStoreID := c.Request().Header.Get("X-Store-Id")
			XChannelID := c.Request().Header.Get("X-Channel-Id")
			XRequestID := c.Request().Header.Get("X-Request-Id")
			XServiceID := c.Request().Header.Get("X-Service-Id")
			XResellerID := c.Request().Header.Get("X-Reseller-Id")
			XIdentity := c.Request().Header.Get("X-Identity")
			XUsername := c.Request().Header.Get("X-Username")

			if storeID == "" {
				c.Request().Header.Set("storeId", XStoreID)
			}
			if channelID == "" {
				c.Request().Header.Set("channelID", XChannelID)
			}
			if requestID == "" {
				c.Request().Header.Set("requestID", XRequestID)
			}
			if serviceID == "" {
				c.Request().Header.Set("serviceID", XServiceID)
			}
			if resellerID == "" {
				c.Request().Header.Set("resellerID", XResellerID)
			}
			if identity == "" {
				c.Request().Header.Set("identity", XIdentity)
			}
			if username == "" {
				c.Request().Header.Set("username", XUsername)
			}

			ctx := logrus.InjectRequestMetadataToContext(c.Request().Context(), c.Request().Header.Get(shared.REQUEST_ID), c.Request().Header.Get(shared.USERNAME), c.Request().URL.Path)
			ctx = context.WithValue(ctx, shared.LOGGER, deps.GetLogger(ctx).WithContext(ctx))
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)
			if err != nil {
				return err
			}

			return nil
		}
	}
}
