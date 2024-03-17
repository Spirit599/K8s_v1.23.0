import numpy as np
import pandas as pd
import json
import sys
from statsmodels.tsa.arima.model import ARIMA
from statsmodels.tools.sm_exceptions import ConvergenceWarning

import warnings

# 忽略特定类型的警告
warnings.filterwarnings("ignore", category=UserWarning)
warnings.filterwarnings("ignore", category=ConvergenceWarning)


# 示例时间序列数据
# data = [1,2,3,4,5,6,7,8,9]
str = sys.argv[1]  # 第一个参数是脚本名称，所以数据是第二个参数

# 解析JSON格式的数据
data = json.loads(str)

# 将数据转换为pandas序列
series = pd.Series(data)

# 定义ARIMA模型
# 这里使用(1,1,1)作为示例参数，实际使用时需要根据数据调整
model = ARIMA(series, order=(1, 1, 1))

# 拟合模型
model_fit = model.fit()

# 打印模型摘要
# print(model_fit.summary())

# 进行预测
preds = model_fit.forecast(steps=1)  # 预测未来1个时间点

print(preds.iloc[0])

# 绘制原始数据和预测数据
# plt.figure(figsize=(10, 6))
# plt.plot(series, label='Original')
# plt.plot(np.arange(len(series), len(series) + 12), preds, label='Forecast', color='red')
# plt.legend()
# plt.show()