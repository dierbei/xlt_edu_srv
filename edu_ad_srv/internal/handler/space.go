package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xlt/edu_srv/edu_ad_srv/global"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/model"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/proto"
)

type SpaceServer struct {
	proto.UnimplementedSpaceServer
}

func (*SpaceServer) GetAllSpaces(ctx context.Context, req *proto.SpaceFilterRequest) (*proto.SpaceListResponse, error) {
	spaces := make([]model.PromotionSpace, 0)
	result := global.MySQLConn.Find(&spaces)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取所有广告位失败")
	}

	SpaceInfoResponse := make([]*proto.SpaceInfoResponse, 0)
	for _, space := range spaces {
		SpaceInfoResponse = append(SpaceInfoResponse, &proto.SpaceInfoResponse{
			Id:       space.ID,
			Name:     space.Name,
			SpaceKey: space.SpaceKey,
			BaseProto: &proto.BaseProto{
				CreateTime: space.CreatedAt.Unix(),
				UpdateTime: space.UpdatedAt.Unix(),
			},
		})
	}

	return &proto.SpaceListResponse{
		Total:     int32(result.RowsAffected),
		SpaceList: SpaceInfoResponse,
	}, nil
}

func (*SpaceServer) GetAdBySpaceKey(ctx context.Context, req *proto.SpaceKeyFilterRequest) (*proto.AdListResponse, error) {
	spaces := make([]model.PromotionSpace, 0)
	result := global.MySQLConn.Where("space_key IN ?", req.SpaceKey).Find(&spaces)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "根据关键词查询广告位失败")
	}

	SpaceAdInfoListResponses := make([]*proto.SpaceAdInfoListResponse, 0)
	for _, space := range spaces {
		sap := &proto.SpaceAdInfoListResponse{
			Id:       space.ID,
			Name:     space.Name,
			SpaceKey: space.SpaceKey,
			BaseProto: &proto.BaseProto{
				CreateTime: space.CreatedAt.Unix(),
				UpdateTime: space.UpdatedAt.Unix(),
			},
		}

		ads := make([]model.PromotionAd, 0)
		if result := global.MySQLConn.Where("space_id = ? and status = ?", space.ID, 1).First(&ads); result.Error != nil {
			return nil, status.Errorf(codes.Internal, "根据广告位查询广告失败")
		}
		AdInfoResponse := make([]*proto.AdInfoResponse, 0)
		for _, ad := range ads {
			adRsp := &proto.AdInfoResponse{
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
				BaseProto: &proto.BaseProto{
					CreateTime: space.CreatedAt.Unix(),
					UpdateTime: space.UpdatedAt.Unix(),
				},
			}
			AdInfoResponse = append(AdInfoResponse, adRsp)
		}
		sap.SpaceAd = AdInfoResponse
		SpaceAdInfoListResponses = append(SpaceAdInfoListResponses, sap)
	}

	return &proto.AdListResponse{SpaceAdInfoListResponses: SpaceAdInfoListResponses}, nil
}
