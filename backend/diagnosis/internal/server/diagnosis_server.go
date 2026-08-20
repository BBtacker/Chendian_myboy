package server

import (
	"context"

	"faceTest/backend/diagnosis/internal/logic"
	"faceTest/backend/diagnosis/internal/svc"
	diagnosispb "faceTest/backend/proto/diagnosis"
)

// DiagnosisServer 诊断服务gRPC服务器
type DiagnosisServer struct {
	svcCtx *svc.ServiceContext
	diagnosispb.UnimplementedDiagnosisServer
}

func NewDiagnosisServer(svcCtx *svc.ServiceContext) *DiagnosisServer {
	return &DiagnosisServer{svcCtx: svcCtx}
}

func (s *DiagnosisServer) SubmitDiagnosis(ctx context.Context, req *diagnosispb.SubmitDiagnosisReq) (*diagnosispb.SubmitDiagnosisResp, error) {
	l := logic.NewSubmitDiagnosisLogic(ctx, s.svcCtx)
	return l.SubmitDiagnosis(req)
}

func (s *DiagnosisServer) GetDiagnosisResult(ctx context.Context, req *diagnosispb.GetResultReq) (*diagnosispb.DiagnosisResultResp, error) {
	l := logic.NewGetDiagnosisResultLogic(ctx, s.svcCtx)
	return l.GetDiagnosisResult(req)
}

func (s *DiagnosisServer) GetTaskStatus(ctx context.Context, req *diagnosispb.GetTaskStatusReq) (*diagnosispb.GetTaskStatusResp, error) {
	l := logic.NewGetTaskStatusLogic(ctx, s.svcCtx)
	return l.GetTaskStatus(req)
}

func (s *DiagnosisServer) ListDiagnosisResults(ctx context.Context, req *diagnosispb.ListResultsReq) (*diagnosispb.ListResultsResp, error) {
	l := logic.NewListDiagnosisResultsLogic(ctx, s.svcCtx)
	return l.ListDiagnosisResults(req)
}

func (s *DiagnosisServer) DeleteDiagnosisResult(ctx context.Context, req *diagnosispb.DeleteResultReq) (*diagnosispb.DeleteResultResp, error) {
	l := logic.NewDeleteDiagnosisResultLogic(ctx, s.svcCtx)
	return l.DeleteDiagnosisResult(req)
}

func (s *DiagnosisServer) BatchDeleteResults(ctx context.Context, req *diagnosispb.BatchDeleteReq) (*diagnosispb.BatchDeleteResp, error) {
	l := logic.NewBatchDeleteResultsLogic(ctx, s.svcCtx)
	return l.BatchDeleteResults(req)
}

func (s *DiagnosisServer) UpdateDiagnosisResult(ctx context.Context, req *diagnosispb.UpdateResultReq) (*diagnosispb.UpdateResultResp, error) {
	l := logic.NewUpdateDiagnosisResultLogic(ctx, s.svcCtx)
	return l.UpdateDiagnosisResult(req)
}

func (s *DiagnosisServer) GetStatistics(ctx context.Context, req *diagnosispb.StatisticsReq) (*diagnosispb.StatisticsResp, error) {
	l := logic.NewGetStatisticsLogic(ctx, s.svcCtx)
	return l.GetStatistics(req)
}

func (s *DiagnosisServer) ExportExcel(ctx context.Context, req *diagnosispb.ExportReq) (*diagnosispb.ExportResp, error) {
	// 委托给报告服务处理
	return &diagnosispb.ExportResp{Code: 0, Msg: "请使用报告服务导出"}, nil
}

func (s *DiagnosisServer) ExportPDF(ctx context.Context, req *diagnosispb.ExportReq) (*diagnosispb.ExportResp, error) {
	// 委托给报告服务处理
	return &diagnosispb.ExportResp{Code: 0, Msg: "请使用报告服务导出"}, nil
}
