package server

import (
	"context"

	"faceTest/backend/report/internal/logic"
	"faceTest/backend/report/internal/svc"
	reportpb "faceTest/backend/proto/report"
)

// ReportServer 报告生成服务gRPC服务器
type ReportServer struct {
	svcCtx *svc.ServiceContext
	reportpb.UnimplementedReportServer
}

func NewReportServer(svcCtx *svc.ServiceContext) *ReportServer {
	return &ReportServer{svcCtx: svcCtx}
}

func (s *ReportServer) ExportExcel(ctx context.Context, req *reportpb.ExportReq) (*reportpb.ExportResp, error) {
	l := logic.NewExportExcelLogic(ctx, s.svcCtx)
	return l.ExportExcel(req)
}

func (s *ReportServer) ExportPDF(ctx context.Context, req *reportpb.ExportReq) (*reportpb.ExportResp, error) {
	l := logic.NewExportPDFLogic(ctx, s.svcCtx)
	return l.ExportPDF(req)
}

func (s *ReportServer) ExportSinglePDF(ctx context.Context, req *reportpb.SingleExportReq) (*reportpb.ExportResp, error) {
	l := logic.NewExportSinglePDFLogic(ctx, s.svcCtx)
	return l.ExportSinglePDF(req)
}
