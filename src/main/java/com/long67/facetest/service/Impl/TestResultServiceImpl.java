package com.long67.facetest.service.Impl;

import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.long67.facetest.mapper.TestResultMapper;
import com.long67.facetest.pojo.testResult;
import com.long67.facetest.service.TestResultService;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.util.List;

// 娣诲姞iText7鐩稿叧鐨勫鍏?
import com.itextpdf.kernel.pdf.PdfDocument;
import com.itextpdf.kernel.pdf.PdfWriter;
import com.itextpdf.layout.Document;
import com.itextpdf.layout.element.Table;
import com.itextpdf.layout.element.Cell;
import com.itextpdf.layout.element.Paragraph;
import com.itextpdf.kernel.font.PdfFont;
import com.itextpdf.kernel.font.PdfFontFactory;
import com.itextpdf.io.font.PdfEncodings;
import com.itextpdf.layout.properties.TextAlignment;

// 娣诲姞Excel鐩稿叧鐨勫鍏?
import org.apache.poi.ss.usermodel.*;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;

@Service
public class TestResultServiceImpl implements TestResultService {
    
    @Autowired
    private TestResultMapper testResultMapper;
    
    /**
     * 鐢ㄦ埛鏌ョ湅鑷繁鐨勬娴嬬粨鏋滐紙鏀寔鍒嗛〉鍜屾潯浠舵绱級
     * @param userId 鐢ㄦ埛ID
     * @param page 椤电爜
     * @param pageSize 姣忛〉澶у皬
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @return 妫€娴嬬粨鏋滃垎椤典俊鎭?
     */
    @Override
    public PageInfo<testResult> getUserTestResults(Integer userId, Integer page, Integer pageSize,
                                                   LocalDateTime startDate, LocalDateTime endDate,
                                                   Boolean isGlandFace, String level) {
        PageHelper.startPage(page, pageSize);
        List<testResult> results = testResultMapper.getUserTestResults(userId, startDate, endDate, isGlandFace, level);
        return new PageInfo<>(results);
    }
    
    /**
     * 鐢ㄦ埛鍒犻櫎鑷繁鐨勬娴嬬粨鏋?
     * @param resultId 妫€娴嬬粨鏋淚D
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎鏄惁鎴愬姛
     */
    @Override
    public boolean deleteUserTestResult(Integer resultId, Integer userId) {
        return testResultMapper.deleteUserTestResult(resultId, userId) > 0;
    }
    
    /**
     * 鐢ㄦ埛鎵归噺鍒犻櫎妫€娴嬬粨鏋?
     * @param ids 妫€娴嬬粨鏋淚D鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎鐨勮褰曟暟
     */
    @Override
    public int deleteBatchUserTestResults(List<Integer> ids, Integer userId) {
        // 鍙傛暟妫€鏌?
        if (ids == null || ids.isEmpty()) {
            return 0;
        }
        
        // 杩囨护鎺夌┖鍊?
        ids.removeIf(id -> id == null);
        if (ids.isEmpty()) {
            return 0;
        }
        
        return testResultMapper.deleteBatchUserTestResults(ids, userId);
    }
    
    /**
     * 瀵煎嚭鐢ㄦ埛妫€娴嬬粨鏋滃埌Excel
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @throws Exception 瀵煎嚭寮傚父
     */
    @Override
    public void exportTestResultsToExcel(HttpServletResponse response, Integer userId, LocalDateTime startDate, LocalDateTime endDate,
                                       Boolean isGlandFace, String level) throws Exception {
        // 鏌ヨ鏁版嵁
        List<testResult> results = testResultMapper.getTestResultsForExport(userId, startDate, endDate, isGlandFace, level);

        // 鍒涘缓宸ヤ綔绨垮拰宸ヤ綔琛?
        Workbook workbook = new XSSFWorkbook();
        Sheet sheet = workbook.createSheet("妫€娴嬭褰?);

        // 鍒涘缓琛ㄥご
        Row headerRow = sheet.createRow(0);
        String[] headers = {"ID", "鐢ㄦ埛ID", "鍥剧墖璺緞", "鏄惁涓鸿吅浣撻潰瀹?, "绛夌骇", "缃俊搴?, "鍙鍖栨弿杩?, "妫€娴嬫椂闂?};
        CellStyle headerStyle = workbook.createCellStyle();
        Font font = workbook.createFont();
        font.setBold(true);
        headerStyle.setFont(font);

        for (int i = 0; i < headers.length; i++) {
            org.apache.poi.ss.usermodel.Cell cell = headerRow.createCell(i);
            cell.setCellValue(headers[i]);
            cell.setCellStyle(headerStyle);
        }

        // 濉厖鏁版嵁
        for (int i = 0; i < results.size(); i++) {
            testResult result = results.get(i);
            Row row = sheet.createRow(i + 1);
            org.apache.poi.ss.usermodel.Cell cell0 = row.createCell(0);
            cell0.setCellValue(result.getId());
            
            org.apache.poi.ss.usermodel.Cell cell1 = row.createCell(1);
            cell1.setCellValue(result.getUserId());
            
            org.apache.poi.ss.usermodel.Cell cell2 = row.createCell(2);
            cell2.setCellValue(result.getImagePath());
            
            org.apache.poi.ss.usermodel.Cell cell3 = row.createCell(3);
            cell3.setCellValue(result.getIsGlandFace() != null ? (result.getIsGlandFace() ? "鏄? : "鍚?) : "");
            
            org.apache.poi.ss.usermodel.Cell cell4 = row.createCell(4);
            cell4.setCellValue(result.getLevel());
            
            org.apache.poi.ss.usermodel.Cell cell5 = row.createCell(5);
            cell5.setCellValue(result.getConfidence() != null ? result.getConfidence() : 0);
            
            org.apache.poi.ss.usermodel.Cell cell6 = row.createCell(6);
            cell6.setCellValue(result.getVisualizationDescription());
            
            org.apache.poi.ss.usermodel.Cell cell7 = row.createCell(7);
            cell7.setCellValue(result.getCreateTime() != null ? result.getCreateTime().toString() : "");
        }

        // 璁剧疆鍝嶅簲澶?
        response.setContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet");
        response.setCharacterEncoding("utf-8");
        String fileName = URLEncoder.encode("妫€娴嬭褰?xlsx", StandardCharsets.UTF_8.toString());
        response.setHeader("Content-Disposition", "attachment; filename=" + fileName);

        // 鍐欏叆鍝嶅簲
        workbook.write(response.getOutputStream());
        workbook.close();
    }
    
    /**
     * 瀵煎嚭鐢ㄦ埛妫€娴嬬粨鏋滃埌PDF
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @throws Exception 瀵煎嚭寮傚父
     */
    @Override
    public void exportTestResultsToPDF(HttpServletResponse response, Integer userId, LocalDateTime startDate, LocalDateTime endDate,
                                      Boolean isGlandFace, String level) throws Exception {
        // 鏌ヨ鏁版嵁
        List<testResult> results = testResultMapper.getTestResultsForExport(userId, startDate, endDate, isGlandFace, level);
        
        // 鐢熸垚HTML鍐呭
        StringBuilder htmlContent = new StringBuilder();
        htmlContent.append("<!DOCTYPE html>");
        htmlContent.append("<html>");
        htmlContent.append("<head>");
        htmlContent.append("<meta charset='UTF-8'>");
        htmlContent.append("<title>妫€娴嬭褰曟姤鍛?/title>");
        htmlContent.append("<style>");
        htmlContent.append("body { font-family: Arial, sans-serif; margin: 20px; }");
        htmlContent.append("h1 { text-align: center; color: #333; }");
        htmlContent.append("table { width: 100%; border-collapse: collapse; margin-top: 20px; }");
        htmlContent.append("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }");
        htmlContent.append("th { background-color: #f2f2f2; font-weight: bold; }");
        htmlContent.append("tr:nth-child(even) { background-color: #f9f9f9; }");
        htmlContent.append(".footer { margin-top: 30px; text-align: center; color: #666; font-size: 12px; }");
        htmlContent.append(".image-preview { max-width: 100px; max-height: 100px; }");
        htmlContent.append("</style>");
        htmlContent.append("</head>");
        htmlContent.append("<body>");
        htmlContent.append("<h1>妫€娴嬭褰曟姤鍛?/h1>");
        htmlContent.append("<table>");
        htmlContent.append("<thead>");
        htmlContent.append("<tr>");
        htmlContent.append("<th>ID</th>");
        htmlContent.append("<th>鐢ㄦ埛ID</th>");
        htmlContent.append("<th>鍥剧墖</th>");
        htmlContent.append("<th>鏄惁涓鸿吅浣撻潰瀹?/th>");
        htmlContent.append("<th>绛夌骇</th>");
        htmlContent.append("<th>缃俊搴?/th>");
        htmlContent.append("<th>鍙鍖栨弿杩?/th>");
        htmlContent.append("<th>妫€娴嬫椂闂?/th>");
        htmlContent.append("</tr>");
        htmlContent.append("</thead>");
        htmlContent.append("<tbody>");
        
        // 濉厖鏁版嵁
        for (testResult result : results) {
            htmlContent.append("<tr>");
            htmlContent.append("<td>").append(result.getId()).append("</td>");
            htmlContent.append("<td>").append(result.getUserId()).append("</td>");
            
            // 澶勭悊鍥剧墖璺緞锛屾坊鍔犲浘鐗囬瑙?
            String imagePath = result.getImagePath();
            if (imagePath != null && !imagePath.isEmpty()) {
                // 濡傛灉鏄浉瀵硅矾寰勶紝娣诲姞鍩虹URL
                if (!imagePath.startsWith("http")) {
                    imagePath = "https://java-web-ai388.oss-cn-beijing.aliyuncs.com/" + imagePath;
                }
                htmlContent.append("<td><a href='").append(imagePath).append("' target='_blank'><img src='").append(imagePath).append("' class='image-preview' alt='妫€娴嬪浘鐗?/></a></td>");
            } else {
                htmlContent.append("<td></td>");
            }
            
            htmlContent.append("<td>").append(result.getIsGlandFace() != null ? (result.getIsGlandFace() ? "鏄? : "鍚?) : "").append("</td>");
            htmlContent.append("<td>").append(result.getLevel()).append("</td>");
            htmlContent.append("<td>").append(result.getConfidence() != null ? result.getConfidence() : 0).append("</td>");
            htmlContent.append("<td>").append(result.getVisualizationDescription() != null ? result.getVisualizationDescription() : "").append("</td>");
            htmlContent.append("<td>").append(result.getCreateTime() != null ? result.getCreateTime().toString() : "").append("</td>");
            htmlContent.append("</tr>");
        }
        
        htmlContent.append("</tbody>");
        htmlContent.append("</table>");
        htmlContent.append("<div class='footer'>");
        htmlContent.append("鎶ュ憡鐢熸垚鏃堕棿: ").append(LocalDateTime.now().toString());
        htmlContent.append("</div>");
        htmlContent.append("</body>");
        htmlContent.append("</html>");
        
        // 璁剧疆鍝嶅簲澶?
        response.setContentType("text/html;charset=utf-8");
        response.setCharacterEncoding("utf-8");
        String fileName = URLEncoder.encode("妫€娴嬭褰?pdf.html", StandardCharsets.UTF_8.toString());
        response.setHeader("Content-Disposition", "attachment; filename=" + fileName);
        
        // 杈撳嚭HTML鍐呭锛岀敤鎴峰彲浠ラ€氳繃娴忚鍣ㄥ彟瀛樹负PDF
        response.getWriter().write(htmlContent.toString());
    }
    
    /**
     * 鐢ㄦ埛鏇存柊鑷繁鐨勬娴嬬粨鏋?
     * @param testResult 妫€娴嬬粨鏋滃璞?
     * @param userId 鐢ㄦ埛ID
     * @return 鏇存柊鏄惁鎴愬姛
     */
    @Override
    public boolean updateUserTestResult(testResult testResult, Integer userId) {
        return testResultMapper.updateUserTestResult(testResult, userId) > 0;
    }
}