package handler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xlt/edu_srv/edu_ad_srv/global"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/model"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/proto"
)

type SpaceServer struct {
	proto.UnimplementedSpaceServer
}

func (*SpaceServer) AdUpdateStatus(ctx context.Context, req *proto.AdUpdateStatusRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdUpdateStatus not implemented")
}

func (*SpaceServer) SpaceSaveOrUpdate(ctx context.Context, req *proto.SpaceSaveOrUpdateRequest) (*empty.Empty, error) {
	space := model.PromotionSpace{}
	global.MySQLConn.Where("id = ?", req.Id).First(&space)
	if req.Name != "" {
		space.Name = req.Name
	}
	if req.SpaceKey != "" {
		space.SpaceKey = req.SpaceKey
	}
	if result := global.MySQLConn.Save(&space); result.Error != nil {
		return nil, status.Errorf(codes.Unimplemented, "新增或者修改广告位失败")
	}
	return &empty.Empty{}, nil
}

func (*SpaceServer) GetSpaceById(ctx context.Context, req *proto.GetSpaceByIdRequest) (*proto.GetSpaceByIdResponse, error) {
	space := model.PromotionSpace{}
	if result := global.MySQLConn.Where("id = ?", req.Id).First(&space); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "根据Id获取广告位失败")
	}
	return &proto.GetSpaceByIdResponse{
		Id:       space.ID,
		Name:     space.Name,
		SpaceKey: space.SpaceKey,
		BaseProto: &proto.BaseProto{
			CreateTime: space.CreatedAt.Unix(),
			UpdateTime: space.UpdatedAt.Unix(),
		},
	}, nil
}

func (*SpaceServer) GetAllSpaces(ctx context.Context, req *proto.GetAllSpacesRequest) (*proto.SpaceInfoListResponse, error) {
	spaces := make([]model.PromotionSpace, 0)
	if result := global.MySQLConn.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&spaces); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取所有的广告位失败")
	}

	spaceInfoResponseList := make([]*proto.SpaceInfoResponse, 0)
	for _, space := range spaces {
		spaceInfoResponse := &proto.SpaceInfoResponse{
			Id:       space.ID,
			Name:     space.Name,
			SpaceKey: space.SpaceKey,
			BaseProto: &proto.BaseProto{
				CreateTime: space.CreatedAt.Unix(),
				UpdateTime: space.UpdatedAt.Unix(),
			},
		}
		spaceInfoResponseList = append(spaceInfoResponseList, spaceInfoResponse)
	}

	return &proto.SpaceInfoListResponse{SpaceInfoResponse: spaceInfoResponseList}, nil
}

func (*SpaceServer) AdSaveOrUpdate(ctx context.Context, req *proto.AdSaveOrUpdateRequest) (*empty.Empty, error) {
	ad := model.PromotionAd{}
	global.MySQLConn.Where("id = ?", req.Id).First(&ad)

	if req.SpaceId != 0 {
		if result := global.MySQLConn.Where("id = ?", req.SpaceId).First(&model.PromotionSpace{}); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.Internal, "广告位不存在")
		}
		ad.SpaceID = int(req.SpaceId)
	}

	if req.Name != "" {
		ad.Name = req.Name
	}

	if ad.Keyword != "" {
		ad.Keyword = req.Keyword
	}
	if ad.HtmlContent != "" {
		ad.HtmlContent = req.HtmlContent
	}
	if ad.Text != "" {
		ad.Text = req.Text
	}
	if ad.Link != "" {
		ad.Link = req.Link
	}
	if ad.StartTime != time.Unix(0, 0) {
		ad.StartTime = time.Unix(req.StartTime, 0)
	}
	if ad.EndTime != time.Unix(0, 0) {
		ad.EndTime = time.Unix(req.EndTime, 0)
	}
	if req.Status != 0 {
		ad.Status = int(req.Status)
	}
	if req.Priority != 0 {
		ad.Priority = int(req.Priority)
	}
	if req.Img != "" {
		ad.Img = req.Img
	}

	if result := global.MySQLConn.Save(&ad); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "新增或者修改广告信息失败")
	}
	return &empty.Empty{}, nil
}

func (*SpaceServer) GetAdById(ctx context.Context, req *proto.GetAdByIdRequest) (*proto.AdInfoResponse, error) {
	ad := model.PromotionAd{}
	result := global.MySQLConn.Where("id = ?", req.Id).First(&ad)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "广告不存在")
	}
	return &proto.AdInfoResponse{
		Name:        ad.Name,
		SpaceId:     int32(ad.SpaceID),
		Keyword:     ad.Keyword,
		HtmlContent: ad.HtmlContent,
		Text:        ad.Text,
		Link:        ad.Link,
		StartTime:   ad.StartTime.Unix(),
		EndTime:     ad.EndTime.Unix(),
		Status:      int32(ad.Status),
		Priority:    int32(ad.Priority),
		Img:         ad.Img,
		Id:          ad.ID,
		BaseProto: &proto.BaseProto{
			CreateTime: ad.BaseModel.CreatedAt.Unix(),
			UpdateTime: ad.BaseModel.UpdatedAt.Unix(),
		},
	}, nil
}

func (*SpaceServer) GetAdList(ctx context.Context, req *proto.GetAdListRequest) (*proto.AdInfoListResponse, error) {
	ads := make([]model.PromotionAd, 0)
	result := global.MySQLConn.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&ads)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取所有广告信息失败")
	}

	adInfoResponseList := make([]*proto.AdInfoResponse, 0)
	for _, ad := range ads {
		adInfoResponse := &proto.AdInfoResponse{
			Name:        ad.Name,
			SpaceId:     int32(ad.SpaceID),
			Keyword:     ad.Keyword,
			HtmlContent: ad.HtmlContent,
			Text:        ad.Text,
			Link:        ad.Link,
			StartTime:   ad.StartTime.Unix(),
			EndTime:     ad.EndTime.Unix(),
			Status:      int32(ad.Status),
			Priority:    int32(ad.Priority),
			Img:         ad.Img,
			Id:          ad.ID,
			BaseProto: &proto.BaseProto{
				CreateTime: ad.CreatedAt.Unix(),
				UpdateTime: ad.UpdatedAt.Unix(),
			},
		}
		adInfoResponseList = append(adInfoResponseList, adInfoResponse)
	}

	return &proto.AdInfoListResponse{AdInfoResponses: adInfoResponseList}, nil
}

func (*SpaceServer) GetAllAds(ctx context.Context, req *proto.GetAllAdsRequest) (*proto.GetAllAdsResponse, error) {
	space := model.PromotionSpace{}
	if result := global.MySQLConn.Where("space_key = ?", req.SpaceKey).First(&space); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "根据关键词查询广告位失败")
	}

	ads := make([]model.PromotionAd, 0)
	if result := global.MySQLConn.Where("space_id = ?", space.ID).Find(&ads); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "根据广告位ID查询广告查询失败")
	}

	adInfoResponseList := make([]*proto.AdInfoResponse, 0)
	for _, ad := range ads {
		adInfoResponse := &proto.AdInfoResponse{
			Name:        ad.Name,
			SpaceId:     int32(ad.SpaceID),
			Keyword:     ad.Keyword,
			HtmlContent: ad.HtmlContent,
			Text:        ad.Text,
			Link:        ad.Link,
			StartTime:   ad.StartTime.Unix(),
			EndTime:     ad.EndTime.Unix(),
			Status:      int32(ad.Status),
			Priority:    int32(ad.Priority),
			Img:         ad.Img,
			Id:          ad.ID,
			BaseProto: &proto.BaseProto{
				CreateTime: ad.CreatedAt.Unix(),
				UpdateTime: ad.UpdatedAt.Unix(),
			},
		}
		adInfoResponseList = append(adInfoResponseList, adInfoResponse)
	}

	return &proto.GetAllAdsResponse{
		Id:       space.ID,
		Name:     space.Name,
		SpaceKey: space.SpaceKey,
		BaseProto: &proto.BaseProto{
			CreateTime: space.CreatedAt.Unix(),
			UpdateTime: space.UpdatedAt.Unix(),
		},
		AdInfoListRsp: adInfoResponseList,
	}, nil
}

//func (*SpaceServer) GetAllSpaces(ctx context.Context, req *proto.SpaceFilterRequest) (*proto.SpaceListResponse, error) {
//	spaces := make([]model.PromotionSpace, 0)
//	result := global.MySQLConn.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&spaces)
//	if result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "获取所有广告位失败")
//	}
//
//	SpaceInfoResponse := make([]*proto.SpaceInfoResponse, 0)
//	for _, space := range spaces {
//		SpaceInfoResponse = append(SpaceInfoResponse, &proto.SpaceInfoResponse{
//			Id:       space.ID,
//			Name:     space.Name,
//			SpaceKey: space.SpaceKey,
//			BaseProto: &proto.BaseProto{
//				CreateTime: space.CreatedAt.Unix(),
//				UpdateTime: space.UpdatedAt.Unix(),
//			},
//		})
//	}
//
//	return &proto.SpaceListResponse{
//		Total:     int32(result.RowsAffected),
//		SpaceList: SpaceInfoResponse,
//	}, nil
//}
//
//func (*SpaceServer) GetAdBySpaceKey(ctx context.Context, req *proto.SpaceKeyFilterRequest) (*proto.SpaceAdInfoListResponse, error) {
//	space := model.PromotionSpace{}
//	result := global.MySQLConn.Where("space_key = ?", req.SpaceKey).Find(&space)
//	if result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "根据关键词查询广告位失败")
//	}
//
//	SpaceAdInfoListResponse := &proto.SpaceAdInfoListResponse{
//		Id:       space.ID,
//		Name:     space.Name,
//		SpaceKey: space.SpaceKey,
//		BaseProto: &proto.BaseProto{
//			CreateTime: space.CreatedAt.Unix(),
//			UpdateTime: space.UpdatedAt.Unix(),
//		},
//	}
//
//	ads := make([]model.PromotionAd, 0)
//	if result := global.MySQLConn.Where("space_id = ? and status = ?", space.ID, 1).Find(&ads); result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "根据广告位查询广告失败")
//	}
//	AdInfoResponse := make([]*proto.AdInfoResponse, 0)
//	for _, ad := range ads {
//		adRsp := &proto.AdInfoResponse{
//			Id:          ad.ID,
//			Name:        ad.Name,
//			SpaceId:     int32(ad.SpaceID),
//			Keyword:     ad.Keyword,
//			HtmlContent: ad.HtmlContent,
//			Text:        ad.Text,
//			Link:        ad.Link,
//			StartTime:   ad.StartTime.Unix(),
//			EndTime:     ad.EndTime.Unix(),
//			Status:      int32(ad.Status),
//			Priority:    int32(ad.Priority),
//			Img:         ad.Img,
//			BaseProto: &proto.BaseProto{
//				CreateTime: space.CreatedAt.Unix(),
//				UpdateTime: space.UpdatedAt.Unix(),
//			},
//		}
//		AdInfoResponse = append(AdInfoResponse, adRsp)
//	}
//	SpaceAdInfoListResponse.SpaceAd = AdInfoResponse
//
//	return SpaceAdInfoListResponse, nil
//}
//
//func (*SpaceServer) SaveOrUpdateSpace(ctx context.Context, req *proto.SpaceInfoRequest) (*empty.Empty, error) {
//	space := model.PromotionSpace{}
//	if result := global.MySQLConn.Where("id = ?", req.Id).First(&space); result.RowsAffected == 0 {
//		space := model.PromotionSpace{
//			Name:     req.Name,
//			SpaceKey: req.SpaceKey,
//		}
//		if result := global.MySQLConn.Create(&space); result.Error != nil {
//			return nil, status.Errorf(codes.Internal, "创建广告信息失败")
//		}
//	} else {
//		space := model.PromotionSpace{
//			BaseModel: model.BaseModel{
//				ID:        req.Id,
//				CreatedAt: space.CreatedAt,
//			},
//			Name:     req.Name,
//			SpaceKey: req.SpaceKey,
//		}
//		if result := global.MySQLConn.Save(&space); result.Error != nil {
//			return nil, status.Errorf(codes.Internal, "更新广告信息失败")
//		}
//	}
//	return &empty.Empty{}, nil
//}
//
//func (*SpaceServer) GetSpaceById(ctx context.Context, req *proto.SpaceByIdRequest) (*proto.SpaceInfoResponse, error) {
//	space := model.PromotionSpace{}
//	if result := global.MySQLConn.Where("id = ?", req.Id).First(&space); result.RowsAffected == 0 {
//		return nil, status.Errorf(codes.Internal, "广告位不存在")
//	}
//	return &proto.SpaceInfoResponse{
//		Id:       space.ID,
//		Name:     space.Name,
//		SpaceKey: space.SpaceKey,
//		BaseProto: &proto.BaseProto{
//			CreateTime: space.CreatedAt.Unix(),
//			UpdateTime: space.UpdatedAt.Unix(),
//		},
//	}, nil
//}
//
//func (*SpaceServer) GetAllAds(ctx context.Context, req *proto.AdPageRequest) (*proto.SpaceAdInfoListResponse, error) {
//	space := model.PromotionSpace{}
//	if result := global.MySQLConn.First(&space); result.RowsAffected == 0 {
//		return nil, status.Errorf(codes.Internal, "查询广告位信息失败")
//	}
//
//	SpaceAdInfoListResponse := &proto.SpaceAdInfoListResponse{
//		Id:       space.ID,
//		Name:     space.Name,
//		SpaceKey: space.SpaceKey,
//		BaseProto: &proto.BaseProto{
//			CreateTime: space.CreatedAt.Unix(),
//			UpdateTime: space.UpdatedAt.Unix(),
//		},
//	}
//
//	ads := make([]model.PromotionAd, 0)
//	result := global.MySQLConn.Where("space_id = ?", space.ID).Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&ads)
//	if result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "查询所有广告信息失败")
//	}
//
//	AdInfoResponses := make([]*proto.AdInfoResponse, 0)
//	for _, ad := range ads {
//		adinfoRsp := &proto.AdInfoResponse{
//			Id:          ad.ID,
//			Name:        ad.Name,
//			SpaceId:     int32(ad.SpaceID),
//			Keyword:     ad.Keyword,
//			HtmlContent: ad.HtmlContent,
//			Text:        ad.Text,
//			Link:        ad.Link,
//			StartTime:   ad.StartTime.Unix(),
//			EndTime:     ad.EndTime.Unix(),
//			Status:      int32(ad.Status),
//			Priority:    int32(ad.Priority),
//			Img:         ad.Img,
//			BaseProto: &proto.BaseProto{
//				CreateTime: ad.CreatedAt.Unix(),
//				UpdateTime: ad.UpdatedAt.Unix(),
//			},
//		}
//		AdInfoResponses = append(AdInfoResponses, adinfoRsp)
//	}
//
//	SpaceAdInfoListResponse.SpaceAd = AdInfoResponses
//
//	return SpaceAdInfoListResponse, nil
//}
//
//func (*SpaceServer) GetAllAds1(ctx context.Context, req *proto.AdPageRequest) (*proto.AdListResponse, error) {
//	spaces := make([]model.PromotionSpace, 0)
//	if result := global.MySQLConn.Find(&spaces); result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "查询所有广告位信息失败")
//	}
//
//	adListResponse := &proto.AdListResponse{}
//	for _, space := range spaces {
//		SpaceAdInfoListResponse := &proto.SpaceAdInfoListResponse{
//			Id:       space.ID,
//			Name:     space.Name,
//			SpaceKey: space.SpaceKey,
//			BaseProto: &proto.BaseProto{
//				CreateTime: space.CreatedAt.Unix(),
//				UpdateTime: space.UpdatedAt.Unix(),
//			},
//		}
//
//		ads := make([]model.PromotionAd, 0)
//		result := global.MySQLConn.Where("space_id = ?", space.ID).Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&ads)
//		if result.Error != nil {
//			return nil, status.Errorf(codes.Internal, "查询所有广告信息失败")
//		}
//
//		AdInfoResponses := make([]*proto.AdInfoResponse, 0)
//		for _, ad := range ads {
//			adinfoRsp := &proto.AdInfoResponse{
//				Id:          ad.ID,
//				Name:        ad.Name,
//				SpaceId:     int32(ad.SpaceID),
//				Keyword:     ad.Keyword,
//				HtmlContent: ad.HtmlContent,
//				Text:        ad.Text,
//				Link:        ad.Link,
//				StartTime:   ad.StartTime.Unix(),
//				EndTime:     ad.EndTime.Unix(),
//				Status:      int32(ad.Status),
//				Priority:    int32(ad.Priority),
//				Img:         ad.Img,
//				BaseProto: &proto.BaseProto{
//					CreateTime: ad.CreatedAt.Unix(),
//					UpdateTime: ad.UpdatedAt.Unix(),
//				},
//			}
//			AdInfoResponses = append(AdInfoResponses, adinfoRsp)
//		}
//
//		SpaceAdInfoListResponse.SpaceAd = AdInfoResponses
//
//		adListResponse.SpaceAdInfoListResponses = append(adListResponse.SpaceAdInfoListResponses, SpaceAdInfoListResponse)
//	}
//
//	return adListResponse, nil
//}
//
//func (*SpaceServer) SaveOrUpdateAd(ctx context.Context, req *proto.AdInfoRequest) (*empty.Empty, error) {
//	ad := model.PromotionAd{}
//	global.MySQLConn.Where("id = ?", req.Id).First(&ad)
//	ad.Name = req.Name
//	ad.SpaceID = int(req.SpaceId)
//	ad.Keyword = req.Keyword
//	ad.HtmlContent = req.HtmlContent
//	ad.Text = req.Text
//	ad.Link = req.Link
//	ad.StartTime = time.Unix(req.StartTime, 0)
//	ad.EndTime = time.Unix(req.EndTime, 0)
//	ad.Status = int(req.Status)
//	ad.Priority = int(req.Priority)
//	ad.Img = req.Img
//	if result := global.MySQLConn.Save(&ad); result.Error != nil {
//		return nil, status.Errorf(codes.Internal, "创建广告失败")
//	}
//	return &empty.Empty{}, nil
//}
//
//func (*SpaceServer) GetAdById(ctx context.Context, req *proto.AdByIdRequest) (*proto.AdInfoResponse, error) {
//	ad := model.PromotionAd{}
//	result := global.MySQLConn.Where("id = ?", req.Id).First(&ad)
//	if result.RowsAffected == 0 {
//		return nil, status.Errorf(codes.Internal, "广告不存在")
//	}
//	return &proto.AdInfoResponse{
//		Name:        ad.Name,
//		SpaceId:     int32(ad.SpaceID),
//		Keyword:     ad.Keyword,
//		HtmlContent: ad.HtmlContent,
//		Text:        ad.Text,
//		Link:        ad.Link,
//		StartTime:   ad.StartTime.Unix(),
//		EndTime:     ad.EndTime.Unix(),
//		Status:      int32(ad.Status),
//		Priority:    int32(ad.Priority),
//		Img:         ad.Img,
//		Id:          ad.ID,
//		BaseProto: &proto.BaseProto{
//			CreateTime: ad.CreatedAt.Unix(),
//			UpdateTime: ad.UpdatedAt.Unix(),
//		},
//	}, nil
//}
