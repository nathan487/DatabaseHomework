# SQL查询类型使用情况统计

## 一、已使用的查询类型

### 1. 单表查询 ✅
**位置**: [service/activity_service.go](service/activity_service.go), [service/user_service.go](service/user_service.go)

**示例代码**:
```go
// ListActivities - 单表查询 + WHERE + ORDER BY
query := config.DB.Model(&model.Activity{})
query = query.Where("status = ?", "active")
if deptID != nil {
    query = query.Where("dept_id = ?", *deptID)
}
if err := query.Order("activity_time desc").Find(&activities).Error; err != nil { ... }

// SearchActivities - 单表查询 + LIKE + WHERE
config.DB.Where("status = ? AND title LIKE ?", "active", "%"+keyword+"%")
```

**使用场景**:
- 查询所有活动列表（按部门、分类筛选）
- 搜索活动（关键词搜索）
- 查询用户信息
- 查询角色信息

---

### 2. JOIN查询（内连接） ✅
**位置**: [service/application_service.go](service/application_service.go)

**示例代码**:
```go
// GetActivityApplications - 2表INNER JOIN
config.DB.Select("Application.application_id, Application.user_id, User.username, Application.apply_time, Application.current_status").
    Joins("JOIN User ON Application.user_id = User.user_id").
    Where("Application.activity_id = ?", activityID).
    Find(&applications)

// GetUserApplications - 2表INNER JOIN
config.DB.Select("Application.application_id, User.user_id, Activity.activity_id, Activity.title, Application.apply_time, Application.current_status").
    Joins("JOIN Activity ON Application.activity_id = Activity.activity_id").
    Where("Application.user_id = ?", userID).
    Find(&applications)
```

**使用场景**:
- 查询某活动的所有报名者及其信息
- 查询用户报名的所有活动及其状态

---

### 3. 条件查询（WHERE子句） ✅
**位置**: 整个项目的service层文件

**常见条件**:
- `WHERE status = 'active'` - 查询活跃活动
- `WHERE user_id = ?` - 按用户ID查询
- `WHERE activity_id = ?` - 按活动ID查询
- `WHERE username = ?` - 按用户名查询
- `WHERE current_status IN ('pending', 'approved')` - 按报名状态查询

---

### 4. 排序查询（ORDER BY） ✅
**位置**: [service/activity_service.go](service/activity_service.go), [service/schedule_service.go](service/schedule_service.go)

**示例代码**:
```go
.Order("activity_time desc")  // 活动时间倒序
.Order("apply_time desc")      // 报名时间倒序
```

**使用场景**:
- 按活动时间排序显示活动列表
- 按报名时间排序显示报名记录

---

### 5. 子查询（IN子句） ✅
**位置**: [service/activity_service.go](service/activity_service.go)

**示例代码**:
```go
// DeleteActivity - 删除活动前先删除相关的状态日志
config.DB.Delete(&model.ApplicationStatusLog{}, 
    "application_id IN (SELECT application_id FROM Application WHERE activity_id = ?)", 
    activityID)
```

**使用场景**:
- 删除活动时，先删除所有相关的报名记录的状态日志

---

## 二、未使用的查询类型

### 1. 聚合函数 + GROUP BY ✅ 已实现
**位置**: [service/statistics_service.go](service/statistics_service.go)

**实现的函数** (聚合函数+GROUP BY+HAVING):

#### 1. GetDeptStatistics() - 部门统计
```sql
SELECT d.dept_id, d.dept_name,
    COUNT(DISTINCT a.activity_id) as activity_count,
    IFNULL(AVG(a.max_people), 0) as avg_capacity,
    COUNT(DISTINCT a.creator_id) as creator_count
FROM Dept d
LEFT JOIN Activity a ON d.dept_id = a.dept_id
GROUP BY d.dept_id, d.dept_name
ORDER BY activity_count DESC
```
**API路由**: GET `/statistics/departments`
**包含聚合函数**: COUNT(DISTINCT), AVG()

#### 2. GetCategoryStatistics() - 分类统计
```sql
SELECT ac.category_id, ac.category_name,
    COUNT(DISTINCT a.activity_id) as activity_count,
    COALESCE(COUNT(ap.application_id), 0) as total_applications,
    COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count
FROM ActivityCategory ac
LEFT JOIN Activity a ON ac.category_id = a.category_id
LEFT JOIN Application ap ON a.activity_id = ap.activity_id
GROUP BY ac.category_id, ac.category_name
ORDER BY activity_count DESC
```
**API路由**: GET `/statistics/categories`
**包含聚合函数**: COUNT(DISTINCT), SUM(CASE WHEN...)

#### 3. GetUserActivityStatistics() - 用户活跃度统计
```sql
SELECT u.user_id, u.username,
    COUNT(ap.application_id) as total_applied,
    COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
    COALESCE(SUM(CASE WHEN ap.current_status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected_count
FROM User u
LEFT JOIN Application ap ON u.user_id = ap.user_id
GROUP BY u.user_id, u.username
HAVING COUNT(ap.application_id) > 0
ORDER BY approved_count DESC
```
**API路由**: GET `/statistics/users`
**包含聚合函数**: COUNT(), SUM(CASE WHEN...), HAVING

#### 4. GetActivityPopularity() - 活动热度排行
```sql
SELECT a.activity_id, a.title, a.max_people,
    COALESCE(COUNT(ap.application_id), 0) as application_count,
    COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
    ROUND(COALESCE(COUNT(ap.application_id), 0) / a.max_people * 100, 2) as fill_rate
FROM Activity a
LEFT JOIN Application ap ON a.activity_id = ap.activity_id
GROUP BY a.activity_id, a.title, a.max_people
ORDER BY application_count DESC
LIMIT 10
```
**API路由**: GET `/statistics/activities/popularity`
**包含聚合函数**: COUNT(), SUM(CASE WHEN...), ROUND()

#### 5. GetAdminCreationStatistics() - 管理员创建统计
```sql
SELECT u.user_id, u.username,
    COUNT(DISTINCT a.activity_id) as created_activities,
    COALESCE(COUNT(ap.application_id), 0) as total_applications,
    COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_applications
FROM User u
LEFT JOIN Activity a ON u.user_id = a.creator_id
LEFT JOIN Application ap ON a.activity_id = ap.activity_id
WHERE u.role_id = 2
GROUP BY u.user_id, u.username
ORDER BY created_activities DESC
```
**API路由**: GET `/statistics/admins`
**包含聚合函数**: COUNT(DISTINCT), COUNT(), SUM(CASE WHEN...)

**前端集成**: [index.html](index.html) - openStatisticsModal() 函数在系统统计仪表盘中展示所有5个聚合统计

**使用场景**:
- 管理员可在统计仪表盘中查看各部门活动情况
- 按分类统计活动和报名分布
- 查看用户活跃度排名（TOP 10）
- 查看活动热度排行（填充率）
- 查看各管理员的创建和批准情况

---

### 2. 日期时间函数 ✅ 已实现
**位置**: [service/statistics_service.go](service/statistics_service.go)

**实现的函数** (日期时间函数):

#### GetUpcomingActivities() - 即将开始的活动
```sql
SELECT a.activity_id, a.title, a.description, 
    d.dept_name, ac.category_name, a.max_people,
    DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
    DATEDIFF(a.activity_time, NOW()) as days_remaining,
    HOUR(TIMEDIFF(a.activity_time, NOW())) as hours_remaining
FROM Activity a
LEFT JOIN Dept d ON a.dept_id = d.dept_id
LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
WHERE a.activity_time > NOW()
    AND a.activity_time <= DATE_ADD(NOW(), INTERVAL 3 DAY)
    AND a.status = 'active'
ORDER BY a.activity_time ASC
```
**API路由**: GET `/activities/upcoming`
**包含日期函数**: NOW(), DATE_ADD(), DATEDIFF(), TIMEDIFF(), DATE_FORMAT(), HOUR()

**前端集成**: [index.html](index.html) - loadUpcomingActivities() 函数在首页显示"即将开始的活动（3天内）"卡片列表

**使用场景**:
- 用户首页显示3天内即将开始的活动
- 按时间倒计时排序（天数、小时数显示）
- 颜色提醒：1天内红色⚠️ | 1-2天黄色⚠️ | 2天以上绿色✅
- 支持快速报名按钮，提升用户参与度

---

### 3. 关联子查询（EXISTS/NOT EXISTS） ❌
**包括**: DATE(), DATEDIFF(), DATE_SUB(), DATE_FORMAT(), NOW()等

**缺失场景**:
- 查询最近N天的活动
- 统计按月的活动数量
- 计算活动距离现在的天数
- 按年月分组统计

**应用场景示例**:
```sql
-- 查询最近7天的活动
SELECT * FROM Activity WHERE activity_time >= DATE_SUB(NOW(), INTERVAL 7 DAY);

-- 按月统计活动数
SELECT DATE_FORMAT(activity_time, '%Y-%m'), COUNT(*) FROM Activity GROUP BY DATE_FORMAT(activity_time, '%Y-%m');

-- 计算活动相对于现在的天数
SELECT activity_id, DATEDIFF(activity_time, NOW()) as days_from_now FROM Activity;
```

---

### 3. 关联子查询（EXISTS/NOT EXISTS） ✅ 已实现
**位置**: [service/statistics_service.go](service/statistics_service.go)

**实现的函数** (关联子查询):

#### GetPopularActivities() - 热门活动（有批准报名的）
```sql
SELECT a.activity_id, a.title, 
    d.dept_name, ac.category_name, u.username as creator_name,
    a.max_people,
    COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
    COALESCE(SUM(CASE WHEN ap.current_status = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
    COALESCE(SUM(CASE WHEN ap.current_status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected_count,
    ROUND(COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) / a.max_people * 100, 2) as fill_rate
FROM Activity a
LEFT JOIN Dept d ON a.dept_id = d.dept_id
LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
LEFT JOIN User u ON a.creator_id = u.user_id
WHERE EXISTS (
    SELECT 1 FROM Application ap 
    WHERE a.activity_id = ap.activity_id 
    AND ap.current_status = 'approved'
)
GROUP BY a.activity_id, a.title, d.dept_name, ac.category_name, u.username, a.max_people
ORDER BY approved_count DESC
LIMIT 10
```
**API路由**: GET `/activities/popular`
**包含关联子查询**: EXISTS 子句检查是否存在已批准的报名记录

**前端集成**: [index.html](index.html) - openStatisticsModal() 函数在统计仪表盘中展示热门活动排行表

**使用场景**:
- 管理员统计仪表盘显示TOP 10热门活动
- 只显示有已批准报名的活动（使用EXISTS过滤）
- 显示各活动的报名统计（已批准、待审批、已拒绝、填充率）

---

### 4. 多表JOIN（3个及以上表） ✅ 已实现
**位置**: [service/activity_service.go](service/activity_service.go)

**实现的函数** (4表JOIN):

#### GetActivityDetail() - 活动详情（4表JOIN）
```sql
SELECT a.activity_id, a.title, a.description, a.location,
    DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
    a.max_people, a.status,
    DATE_FORMAT(a.create_time, '%Y-%m-%d %H:%i') as create_time,
    a.dept_id, COALESCE(d.dept_name, '') as dept_name,
    a.category_id, COALESCE(ac.category_name, '') as category_name,
    a.creator_id, COALESCE(u.username, '') as creator_name
FROM Activity a
LEFT JOIN Dept d ON a.dept_id = d.dept_id
LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
LEFT JOIN User u ON a.creator_id = u.user_id
WHERE a.activity_id = ?
```
**API路由**: GET `/activities/:id`
**包含多表JOIN**: 4个表（Activity、Dept、ActivityCategory、User）

**前端集成**: [index.html](index.html) - showActivityDetail() 函数显示活动详情

**使用场景**:
- 用户点击"查看详情"时显示完整的活动信息
- 一次查询获取：活动基本信息 + 部门名称 + 分类名称 + 创建者用户名
- 优化显示：显示"志愿服务部 / 社区服务"而不是"部门 1"

---

### 5. 外连接（LEFT JOIN / RIGHT JOIN / FULL OUTER JOIN） ✅
**场景**: 保留一方的所有记录

**缺失场景**:
- 查询所有用户及其报名信息（包括未报名用户）
- 查询所有活动及其报名者（包括无人报名的活动）
- 查询所有部门及其活动（包括无活动的部门）

**应用场景示例**:
```sql
-- 所有用户及其报名情况
SELECT u.user_id, u.username, a.activity_id, a.title, ap.apply_time
FROM User u
LEFT JOIN Application ap ON u.user_id = ap.user_id
LEFT JOIN Activity a ON ap.activity_id = a.activity_id;

-- 所有部门及其活动
SELECT d.dept_id, d.dept_name, COUNT(a.activity_id) as activity_count
FROM Dept d
LEFT JOIN Activity a ON d.dept_id = a.dept_id
GROUP BY d.dept_id, d.dept_name;
```

---

### 6. 集合操作（UNION / NOT IN / NOT EXISTS） ✅ 已实现
**位置**: [service/activity_service.go](service/activity_service.go)

**实现的函数** (NOT IN集合操作):

#### GetAvailableActivities() - 用户可申请的活动（NOT IN集合操作）
```sql
SELECT a.activity_id, a.title, a.description, a.location,
    DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
    a.max_people, COALESCE(COUNT(app.application_id), 0) as current_apply_count,
    (a.max_people - COALESCE(COUNT(app.application_id), 0)) as remaining_slots,
    COALESCE(d.dept_name, '未分配') as dept_name,
    COALESCE(ac.category_name, '未分类') as category_name
FROM Activity a
LEFT JOIN Application app ON a.activity_id = app.activity_id
LEFT JOIN Dept d ON a.dept_id = d.dept_id
LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
WHERE a.status = 'active'
AND a.activity_id NOT IN (
    SELECT DISTINCT activity_id FROM Application WHERE user_id = ?
)
GROUP BY a.activity_id, a.title, a.description, a.location, a.activity_time, 
    a.max_people, d.dept_name, ac.category_name
HAVING remaining_slots > 0
ORDER BY a.activity_time ASC
```
**API路由**: GET `/activities/available`
**包含集合操作**: NOT IN 子句排除已申请过的活动

**前端集成**: [index.html](index.html) - loadAvailableActivities() 函数在用户首页显示"我可以申请的活动"卡片列表

**使用场景**:
- 用户首页显示"我可以申请的活动"
- 排除已申请过的活动（NOT IN集合操作）
- 只显示还有空位的活动（HAVING remaining_slots > 0）
- 显示报名进度条和剩余席位数
- 支持快速报名按钮

---

### 7. 自连接（Self JOIN） ✅ 已实现
**位置**: [service/activity_service.go](service/activity_service.go)

**实现的函数** (自连接):

#### GetAvailableActivities() - 可申请活动排除时间冲突
```sql
-- 完整SQL：自连接部分排除与用户已申请活动时间冲突的活动
AND a.activity_id NOT IN (
    SELECT DISTINCT a2.activity_id
    FROM Activity a2
    JOIN Activity a1 ON ABS(HOUR(TIMEDIFF(a1.activity_time, a2.activity_time))) < 2
    WHERE a1.activity_id IN (
        SELECT DISTINCT activity_id FROM Application WHERE user_id = ?
    ) AND a1.status = 'active' AND a2.status = 'active'
)
```

**解释**:
- `a2` 代表可申请的活动
- `a1` 代表用户已申请的活动
- `JOIN Activity a1` 实现自连接：同一表的两个别名进行关联
- `ABS(HOUR(TIMEDIFF(...))) < 2` 检查时间差是否在2小时以内
- 通过 `NOT IN` 排除所有有冲突的活动

**API路由**: GET `/activities/available?user_id=X`
**包含自连接**: Activity 表自己与自己连接（a1和a2），通过时间差函数检查2小时内的活动冲突

**前端集成**: [index.html](index.html) - loadAvailableActivities() 函数在用户首页显示"我可以申请的活动"

**使用场景**:
- 用户首页显示"我可以申请的活动"
- 排除已申请过的活动（NOT IN集合操作）
- **排除与用户已申请活动时间冲突的活动（自连接 + 时间函数）** ⭐
- 只显示还有空位的活动（HAVING remaining_slots > 0）
- 避免用户选择时间冲突的活动，提升用户体验

---
**场景**: 表与自身连接

**缺失场景**:
- 查询同一部门在同一时间的活动
- 查询参加相同活动的用户对

**应用场景示例**:
```sql
-- 同一部门同一时间的活动
SELECT a1.activity_id, a1.title, a2.activity_id, a2.title
FROM Activity a1
JOIN Activity a2 ON a1.dept_id = a2.dept_id AND DATE(a1.activity_time) = DATE(a2.activity_time)
WHERE a1.activity_id < a2.activity_id;

-- 报名了相同活动的用户对
SELECT a1.user_id, a2.user_id, a1.activity_id
FROM Application a1
JOIN Application a2 ON a1.activity_id = a2.activity_id AND a1.user_id < a2.user_id;
```

---

### 8. 除法查询 ✅ 已实现
**位置**: [service/statistics_service.go](service/statistics_service.go)

**实现的函数** (除法查询):

#### GetOmnipotentVolunteers() - 全能志愿者（参加所有分类活动的用户）
```sql
-- 除法查询：找出参加了所有分类活动的用户
SELECT 
    u.user_id, 
    u.username,
    COUNT(DISTINCT a.category_id) as categories_participated,
    (SELECT COUNT(*) FROM ActivityCategory) as total_categories_count,
    COUNT(DISTINCT app.application_id) as approved_count
FROM User u
JOIN Application app ON u.user_id = app.user_id AND app.current_status = 'approved'
JOIN Activity a ON app.activity_id = a.activity_id
GROUP BY u.user_id, u.username
HAVING COUNT(DISTINCT a.category_id) = (SELECT COUNT(*) FROM ActivityCategory)
ORDER BY approved_count DESC
```

**API路由**: GET `/statistics/omnipotent-volunteers`

**前端集成**: [index.html](index.html) - openStatisticsModal() 函数在统计仪表盘展示全能志愿者

**使用场景**:
- 管理员仪表盘展示全能志愿者排行
- 识别参加所有分类活动的用户
- 用户激励和社区运营

---

## 三、总体统计

| 查询类型 | 使用情况 | 使用数量 | 说明 |
|---------|---------|---------|------|
| 单表查询 | ✅ 已使用 | 多处 | 广泛用于列表、搜索、查询单个记录 |
| WHERE条件 | ✅ 已使用 | 多处 | 几乎所有查询都有条件过滤 |
| ORDER BY | ✅ 已使用 | 多处+ | 用于排序活动和报名列表 |
| 2表JOIN | ✅ 已使用 | 多处 | 用于获取报名者和活动信息 |
| 子查询(IN) | ✅ 已使用 | 1处 | 仅用于删除时的级联操作 |
| 聚合函数+GROUP BY | ✅ 已使用 | 5处 | 管理员统计仪表盘（部门、分类、用户、活动、管理员） |
| 日期时间函数 | ✅ 已使用 | 1处 | 首页显示3天内即将开始的活动 |
| 关联子查询(EXISTS) | ✅ 已使用 | 1处 | 统计仪表盘热门活动（有批准报名的） |
| 多表JOIN(3+) | ✅ 已使用 | 1处 | 活动详情页面（4表JOIN：活动+部门+分类+创建者） |
| 外连接 | ✅ 已使用 | 7处 | 聚合统计、热门活动、活动详情中使用LEFT JOIN保留空记录 |
| 集合操作(NOT IN) | ✅ 已使用 | 1处 | 我可以申请的活动（排除已申请的活动） |
| 自连接 | ✅ 已使用 | 1处 | 可申请活动中排除时间冲突的活动 |
| 除法查询 | ✅ 已使用 | 1处 | 全能志愿者（参加所有分类活动的用户） |

---

## 四、建议改进方向

### 已完成的功能：

1. **管理员统计仪表盘** ✅ (已实现)
   - 使用聚合函数统计各类数据（5个函数）
   - 使用GROUP BY进行分组统计
   - 使用LEFT JOIN保留空记录显示完整数据
   - 使用CASE WHEN进行条件聚合
   - 使用HAVING进行聚合后过滤

2. **即将开始的活动提醒** ✅ (已实现)
   - 使用NOW()、DATE_ADD()获取当前时间和3天后的时间范围
   - 使用DATEDIFF()和TIMEDIFF()计算剩余时间（天数和小时数）
   - 使用DATE_FORMAT()格式化时间显示
   - 按活动时间升序排列
   - 前端显示倒计时和颜色提醒

3. **热门活动排行** ✅ (已实现)
   - 使用EXISTS关联子查询过滤有已批准报名的活动
   - 统计报名情况（已批准、待审批、已拒绝）
   - 计算活动填充率
   - 管理员仪表盘TOP 10显示

4. **活动详情页面** ✅ (已实现)
   - 使用4表LEFT JOIN（Activity、Dept、ActivityCategory、User）
   - 一次查询获取完整的活动信息
   - 前端显示部门名称和分类名称而不是ID
   - 优化用户体验

### 可以在以下功能中使用未使用的查询类型：

1. **高级搜索功能**
   - 使用多表JOIN（3+表）显示完整信息
   - 使用日期范围过滤（活动时间区间）

2. **用户活跃度报告**
   - 使用关联子查询找出特定类型的用户
   - 例：查找参加过志愿服务部所有活动的用户

3. **数据分析**
   - 使用集合操作进行用户分类
   - 使用除法查询找出高级志愿者
   - 使用自连接找出相同时间段的活动冲突
