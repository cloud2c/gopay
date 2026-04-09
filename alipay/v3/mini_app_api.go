package alipay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloud2c/gopay"
)

// 小程序退回开发 alipay.open.mini.version.audited.cancel
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionAuditedCancel(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionAuditedCancelRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionAuditedCancel, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionAuditedCancelRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序灰度上架 alipay.open.mini.version.gray.online
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionGrayOnline(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionGrayOnlineRsp, err error) {
	err = bm.CheckEmptyError("app_version", "gray_strategy")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionGrayOnline, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionGrayOnlineRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序结束灰度 alipay.open.mini.version.gray.cancel
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionGrayCancel(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionGrayCancelRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionGrayCancel, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionGrayCancelRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序上架 alipay.open.mini.version.online
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionOnline(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionOnlineRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionOnline, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionOnlineRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序下架 alipay.open.mini.version.offline
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionOffline(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionOfflineRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionOffline, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionOfflineRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序回滚 alipay.open.mini.version.rollback
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionRollback(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionRollbackRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionRollback, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionRollbackRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序删除版本 alipay.open.mini.version.delete
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionDelete(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionDeleteRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionDelete, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionDeleteRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序提交审核 alipay.open.mini.version.audit.apply
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionAuditApply(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionAuditApplyRsp, err error) {
	err = bm.CheckEmptyError("app_version", "version_desc")
	if err != nil {
		return nil, err
	}
	if bm.GetString("service_email") == gopay.NULL && bm.GetString("service_phone") == gopay.NULL {
		return nil, errors.New("service_email and service_phone are not allowed to be null at the same time")
	}

	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	res, bs, err := a.doProdPostFile(ctx, bm, v3OpenMiniVersionAuditApply, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionAuditApplyRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序基于模板上传版本 alipay.open.mini.version.upload
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionUpload(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionUploadRsp, err error) {
	err = bm.CheckEmptyError("template_id", "app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniVersionUpload, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionUploadRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 查询使用模板的小程序列表 alipay.open.mini.template.usage.query
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniTemplateUsageQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniTemplateUsageQueryRsp, err error) {
	err = bm.CheckEmptyError("template_id")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3OpenMiniTemplateUsageQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniTemplateUsageQueryRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序查询版本构建状态 alipay.open.mini.version.build.query
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionBuildQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionBuildQueryRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3OpenMiniVersionBuildQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionBuildQueryRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序版本详情查询 alipay.open.mini.version.detail.query
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionDetailQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionDetailQueryRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3OpenMiniVersionDetailQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionDetailQueryRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序版本列表查询 alipay.open.mini.version.list.query
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniVersionListQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniVersionListQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3OpenMiniVersionListQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniVersionListQueryRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序生成体验版 alipay.open.mini.experience.create
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniExperienceCreate(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniExperienceCreateRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniExperienceCreate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniExperienceCreateRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序体验版状态查询接口 alipay.open.mini.experience.query
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniExperienceQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniExperienceQueryRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3OpenMiniExperienceQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniExperienceQueryRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 小程序取消体验版 alipay.open.mini.experience.cancel
// StatusCode = 200 is success
func (a *ClientV3) OpenMiniExperienceCancel(ctx context.Context, bm gopay.BodyMap) (aliRsp *OpenMiniExperienceCancelRsp, err error) {
	err = bm.CheckEmptyError("app_version")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3OpenMiniExperienceCancel, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OpenMiniExperienceCancelRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}
