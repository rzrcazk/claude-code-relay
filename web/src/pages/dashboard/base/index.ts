import dayjs from 'dayjs';
import type { EChartsOption } from 'echarts';

import type { TChartColor } from '@/config/color';
import { getChartListColor } from '@/utils/color';

/** 首页 dashboard 小图表（使用真实趋势数据） */
export function constructMiniChart(type: string, trendData: Array<any> = []) {
  // 取最近7天的数据
  const recentData = trendData.slice(-7);
  const dateArray = recentData.map((item) => dayjs(item.date).format('MM-DD'));

  const datasetAxis = {
    xAxis: {
      type: 'category',
      show: false,
      data: dateArray,
    },
    yAxis: {
      show: false,
      type: 'value',
    },
    grid: {
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
    },
  };

  if (type === 'line') {
    const costData = recentData.map((item) => item.cost);
    const lineDataset = {
      ...datasetAxis,
      color: ['#fff'],
      series: [
        {
          data: costData,
          type,
          showSymbol: true,
          symbol: 'circle',
          symbolSize: 3,
          lineStyle: {
            width: 2,
          },
          smooth: true,
        },
      ],
    };
    return lineDataset;
  }

  // bar chart - 使用tokens数据
  const tokensData = recentData.map((item, index) => ({
    value: item.tokens,
    itemStyle: {
      opacity: index >= recentData.length - 2 ? 1 : 0.6,
    },
  }));

  const barDataset = {
    ...datasetAxis,
    color: getChartListColor(),
    series: [
      {
        data: tokensData,
        type,
        barWidth: 9,
      },
    ],
  };
  return barDataset;
}

/**
 *  线性图表数据源（使用真实数据）
 *
 * @export
 * @param {Array} trendData 趋势数据
 * @returns {*} dataSet
 */
export function getLineChartDataSet({
  trendData = [],
  placeholderColor,
  borderColor,
}: { trendData?: Array<any> } & TChartColor) {
  const timeArray = trendData.map((item) => dayjs(item.date).format('MM-DD'));
  const costArray = trendData.map((item) => item.cost.toFixed(2));
  const requestsArray = trendData.map((item) => item.requests);

  const dataSet = {
    color: getChartListColor(),
    tooltip: {
      trigger: 'axis',
      formatter(params: any) {
        let result = `${params[0].axisValue}<br/>`;
        params.forEach((param: any) => {
          const value = param.seriesName === '费用' ? `$${param.value}` : `${param.value}次`;
          result += `${param.marker}${param.seriesName}: ${value}<br/>`;
        });
        return result;
      },
    },
    grid: {
      left: '0',
      right: '20px',
      top: '5px',
      bottom: '36px',
      containLabel: true,
    },
    legend: {
      left: 'center',
      bottom: '0',
      orient: 'horizontal',
      data: ['费用', '请求数'],
      textStyle: {
        fontSize: 12,
        color: placeholderColor,
      },
    },
    xAxis: {
      type: 'category',
      data: timeArray,
      boundaryGap: false,
      axisLabel: {
        color: placeholderColor,
      },
      axisLine: {
        lineStyle: {
          width: 1,
        },
      },
    },
    yAxis: [
      {
        type: 'value',
        name: '费用($)',
        position: 'left',
        axisLabel: {
          color: placeholderColor,
          formatter: '${value}',
        },
        splitLine: {
          lineStyle: {
            color: borderColor,
          },
        },
      },
      {
        type: 'value',
        name: '请求数',
        position: 'right',
        axisLabel: {
          color: placeholderColor,
        },
        splitLine: {
          show: false,
        },
      },
    ],
    series: [
      {
        name: '费用',
        data: costArray,
        type: 'line',
        yAxisIndex: 0,
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 6,
        itemStyle: {
          borderColor,
          borderWidth: 1,
        },
        areaStyle: {
          opacity: 0.1,
        },
      },
      {
        name: '请求数',
        data: requestsArray,
        type: 'line',
        yAxisIndex: 1,
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 6,
        itemStyle: {
          borderColor,
          borderWidth: 1,
        },
      },
    ],
  };
  return dataSet;
}

/**
 * 获取饼图数据（使用真实模型数据）
 *
 * @export
 * @param {Array} modelStats 模型统计数据
 * @returns {*} dataSet
 */
export function getPieChartDataSet({
  modelStats = [],
  textColor,
  placeholderColor,
  containerColor,
}: { modelStats?: Array<any> } & Record<string, string>): EChartsOption {
  const data = modelStats.map((item) => ({
    value: item.cost,
    name: item.model_name.replace(/^claude-/, ''), // 简化模型名称显示
  }));

  const totalCost = modelStats.reduce((sum, item) => sum + item.cost, 0);

  return {
    color: getChartListColor(),
    tooltip: {
      trigger: 'item',
      formatter: '{a}<br/>{b}: ${c} ({d}%)',
    },
    grid: {
      top: '0',
      right: '0',
    },
    legend: {
      selectedMode: false,
      itemWidth: 12,
      itemHeight: 4,
      textStyle: {
        fontSize: 12,
        color: placeholderColor,
      },
      left: 'center',
      bottom: '0',
      orient: 'horizontal',
      formatter(name: string) {
        return name.length > 15 ? `${name.substring(0, 12)}...` : name;
      },
    },
    series: [
      {
        name: '模型费用分布',
        type: 'pie',
        radius: ['48%', '60%'],
        avoidLabelOverlap: true,
        selectedMode: true,
        silent: false,
        itemStyle: {
          borderColor: containerColor,
          borderWidth: 1,
        },
        label: {
          show: true,
          position: 'center',
          formatter(params: any) {
            if (params.dataIndex === 0) {
              const percent = ((params.value / totalCost) * 100).toFixed(1);
              return [`{value|${percent}%}`, `{name|${params.name}}`].join('\n');
            }
            return '';
          },
          rich: {
            value: {
              color: textColor,
              fontSize: 28,
              fontWeight: 'normal',
              lineHeight: 46,
            },
            name: {
              color: '#909399',
              fontSize: 12,
              lineHeight: 14,
            },
          },
        },
        emphasis: {
          scale: true,
          label: {
            show: true,
            formatter(params: any) {
              const percent = ((params.value / totalCost) * 100).toFixed(1);
              return [`{value|${percent}%}`, `{name|${params.name}}`].join('\n');
            },
            rich: {
              value: {
                color: textColor,
                fontSize: 28,
                fontWeight: 'normal',
                lineHeight: 46,
              },
              name: {
                color: '#909399',
                fontSize: 14,
                lineHeight: 14,
              },
            },
          },
        },
        labelLine: {
          show: false,
        },
        data,
      },
    ],
  };
}
