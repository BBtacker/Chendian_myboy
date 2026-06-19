import React, { useState } from 'react';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  CategoryScale,
  LinearScale,
  BarElement,
  Title
} from 'chart.js';
import { Pie, Bar } from 'react-chartjs-2';

// 注册Chart.js组件
ChartJS.register(
  ArcElement,
  Tooltip,
  Legend,
  CategoryScale,
  LinearScale,
  BarElement,
  Title
);

function AnalysisResult({ result }) {
  const [selectedEmotion, setSelectedEmotion] = useState(null);

  // 模拟API返回的结果数据
  const mockResult = {
    image_url: 'http://localhost:5000/static/result.jpg', // 实际应从API获取
    emotion_counts: {
      happy: 25,
      sad: 10,
      angry: 5,
      surprise: 15,
      fear: 5
    },
    emotion_scores: {
      happy: {
        no_happy: 2,
        slight_happy: 8,
        normal_happy: 10,
        broad_happy: 5
      },
      sad: {
        no_happy: 10, // 注意：这里应该是对应sad的等级，为了演示使用了相同的键名
        slight_happy: 0,
        normal_happy: 0,
        broad_happy: 0
      },
      angry: {
        no_happy: 5,
        slight_happy: 0,
        normal_happy: 0,
        broad_happy: 0
      },
      surprise: {
        no_happy: 15,
        slight_happy: 0,
        normal_happy: 0,
        broad_happy: 0
      },
      fear: {
        no_happy: 5,
        slight_happy: 0,
        normal_happy: 0,
        broad_happy: 0
      }
    }
  };

  // 使用实际结果或模拟数据
  const data = result || mockResult;

  // 饼状图数据
  const pieData = {
    labels: Object.keys(data.emotion_counts),
    datasets: [
      {
        data: Object.values(data.emotion_counts),
        backgroundColor: [
          'rgba(75, 192, 192, 0.8)',  // happy - 青绿色
          'rgba(255, 99, 132, 0.8)',   // sad - 红色
          'rgba(255, 206, 86, 0.8)',   // angry - 黄色
          'rgba(153, 102, 255, 0.8)',  // surprise/surprised - 紫色
          'rgba(255, 159, 64, 0.8)',   // fear - 橙色
          'rgba(201, 203, 207, 0.8)',  // neutral/natural - 灰色
          'rgba(154, 205, 50, 0.8)',   // disgust - 黄绿色
          'rgba(0, 128, 128, 0.8)',    // contempt - 深青色
          'rgba(128, 128, 0, 0.8)'     // sleepy - 橄榄色
        ],
        borderColor: [
          'rgba(75, 192, 192, 1)',     // happy - 青绿色
          'rgba(255, 99, 132, 1)',      // sad - 红色
          'rgba(255, 206, 86, 1)',      // angry - 黄色
          'rgba(153, 102, 255, 1)',     // surprise/surprised - 紫色
          'rgba(255, 159, 64, 1)',      // fear - 橙色
          'rgba(201, 203, 207, 1)',     // neutral/natural - 灰色
          'rgba(154, 205, 50, 1)',      // disgust - 黄绿色
          'rgba(0, 128, 128, 1)',       // contempt - 深青色
          'rgba(128, 128, 0, 1)'        // sleepy - 橄榄色
        ],
        borderWidth: 1
      }
    ]
  };

  // 饼状图选项
  const pieOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'right',
      },
      tooltip: {
        callbacks: {
          label: function(context) {
            const label = context.label || '';
            const value = context.raw || 0;
            const total = context.dataset.data.reduce((a, b) => a + b, 0);
            const percentage = Math.round((value / total) * 100);
            return `${label}: ${value} (${percentage}%)`;
          }
        }
      }
    },
    onClick: (event, elements) => {
      if (elements.length > 0) {
        const index = elements[0].index;
        const emotion = Object.keys(data.emotion_counts)[index];
        setSelectedEmotion(emotion);
      }
    }
  };

  // 条形图数据
  const barData = {
    labels: ['无', '轻微', '中等', '强烈'],
    datasets: [
      {
        label: selectedEmotion || '表情等级分布',
        data: selectedEmotion && data.emotion_scores && data.emotion_scores[selectedEmotion] 
          ? Object.values(data.emotion_scores[selectedEmotion]) 
          : [0, 0, 0, 0],
        backgroundColor: 'rgba(54, 162, 235, 0.6)',
        borderColor: 'rgba(54, 162, 235, 1)',
        borderWidth: 1
      }
    ]
  };

  // 条形图选项
  const barOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          precision: 0
        }
      }
    },
    plugins: {
      title: {
        display: true,
        text: selectedEmotion ? `${selectedEmotion}表情等级分布` : '请选择一个表情'
      }
    }
  };

  return (
    <div className="analysis-result">
      {/* 结果图片展示 */}
      <div className="result-image-container">
        <h3>分析结果图片</h3>
        <div className="result-image-wrapper">
          <img 
            src={data.image_url} 
            alt="分析结果" 
            className="result-image"
            onError={(e) => {
              // 图片加载失败时显示占位图
              e.target.src = 'https://via.placeholder.com/800x400?text=Result+Image';
            }}
          />
        </div>
      </div>

      {/* 统计图表 */}
      <div className="charts-container">
        {/* 饼状图 */}
        <div className="chart-section pie-chart">
          <h3>表情分布统计</h3>
          <div className="chart-wrapper">
            <Pie data={pieData} options={pieOptions} />
          </div>
          <p className="chart-hint">点击饼状图查看表情等级分布</p>
        </div>

        {/* 条形图 */}
        <div className={`chart-section bar-chart ${selectedEmotion ? 'show' : 'hide'}`}>
          <h3>表情等级分布</h3>
          <div className="chart-wrapper">
            {selectedEmotion ? (
              <Bar data={barData} options={barOptions} />
            ) : (
              <div className="bar-chart-placeholder">
                <p>请点击左侧饼状图选择一个表情</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* 统计摘要 */}
      <div className="stats-summary">
        <h3>统计摘要</h3>
        <div className="stats-grid">
          {Object.entries(data.emotion_counts).map(([emotion, count]) => (
            <div key={emotion} className="stat-card">
              <div className="stat-emotion">{emotion}</div>
              <div className="stat-count">{count}</div>
              <div className="stat-percentage">
                {Math.round((count / Object.values(data.emotion_counts).reduce((a, b) => a + b, 0)) * 100)}%
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default AnalysisResult;