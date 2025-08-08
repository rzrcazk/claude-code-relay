## 进入开发
### 项目结构[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84)

正如您初始化项目后可以看到，TDesign Starter 的整个项目的目录结构大致如下：

```
.
├── README.md                         # 说明文档
├── index.html                        # 主 html 文件
├── docs
├── mock                              # mock 目录
│     └── index.ts
├── package.json
├── package-lock.json
├── node_modules                      # 项目依赖
├── public
│     └── favicon.ico
├── src                               # 页面代码
├── .env                              # 生产环境变量
├── .env.development                  # 开发环境变量
├── commitlint.config.js              # commintlint 规范
├── tsconfig.json                     # typescript 配置文件
└── vite.config.js                    # vite 配置文件
```

### 页面代码结构[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E9%A1%B5%E9%9D%A2%E4%BB%A3%E7%A0%81%E7%BB%93%E6%9E%84)

如上图所示，`src`目录下是页面代码，大部分情况下，您只需要增删`src`目录下的代码即可。

`src`内的结构大致如下所示，TDesign 推荐您在使用的过程中，遵守既有的目录结构，以规范项目代码的组织结构。

```
src
├── App.vue
├── apis                              # 请求层
├── assets                            # 资源层
├── components                        # 公共组件层
├── config                            # 配置层
│     ├── global.ts                     # 全局常量配置
│     ├── color.ts                      # 全局主题色彩配置
│     └── style.ts                      # 布局样式配置
├── constants                         # 常量层
│     └── index.ts
├── hooks                             # 钩子层
│     └── index.ts
├── layouts                           # 布局层 可动态调整
│     ├── components                    # 布局组件
│     │     ├── Breadcrumb.vue            # 面包屑组件
│     │     ├── ...
│     │     └── SideNav.vue               # 侧边栏组件
│     ├── frame                         # 嵌入式组件
│     │     └── index.vue
│     ├── setting.vue                   # 配置生成组件
│     ├── blank.vue                     # 空白路由
│     └── index.vue
├── pages                             # 业务模块层
│     ├── dashboard                     # 一个页面组件
│     │     └── base
│     │           ├── components          # 该页面组件用到的子组件
│     │           ├── constants.ts        # 该页面组件用到的常量
│     │           ├── index.ts
│     │           └── index.vue
│     ├── ...
│     └── user
│           ├── constants.ts
│           ├── index.less
│           ├── index.ts
│           └── index.vue
├── router                            # 路由层
├── store                             # Pinia 数据层
│     ├── index.ts
│     └── modules
│           ├── notification.ts
│           ├── ...
│           ├── setting.ts
│           └── user.ts
├── style                             # 样式目录
│     ├── font-family.less              # 字体文件（腾讯体W7）
│     ├── layout.less                   # 全局样式布局
│     ├── reset.less                    # 对默认样式的重置
│     └── variables.less                # 模板样式 token
├── types                             # 类型文件目录
└── utils                             # 工具层
│     ├── route                         # 路由工具封装
│     ├── charts.ts                     # 图表工具封装
│     ├── color.ts                      # 色彩工具封装
│     └── request                       # 请求工具封装
├── permission.ts                     # 权限逻辑
└── main.ts                           # 入口逻辑文件

```

### 环境变量[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)

在项目的根目录，有 `.env` 配置文件，项目会根据启动的命令中的 `mode` 参数，加载指定的配置文件的配置来运行， 如本地环境执行 `npm run dev`，因为对于命令中的`mode` 参数为`development`，项目运行会加载`.env.development`的配置来运行。 项目初始化内置了 `.env.development`、`.env.test` 和 `.env` 分别对应本地开发环境、测试环境 和 生产（正式）环境，也可以根据实际需求继续扩展。

#### 内置的环境变量[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E5%86%85%E7%BD%AE%E7%9A%84%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)

- `VITE_BASE_URL`：项目启动运行默认的 URL
- `VITE_IS_REQUEST_PROXY`： 项目是否启动请求代理
- `VITE_API_URL`: 项目默认请求的 URL
- `VITE_API_URL_PREFIX`：项目默认请求的前缀

### 开始开发[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E5%BC%80%E5%A7%8B%E5%BC%80%E5%8F%91)

#### 新增页面[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E6%96%B0%E5%A2%9E%E9%A1%B5%E9%9D%A2)

在已有 TDesign Starter 项目中，新增页面是非常方便的。

首先，在 `./src/pages` 下，创建新页面的目录以及相关的文件。

```
cd src/pages && mkdir my-new-page

cd my-new-page && touch index.vue  # 可根据实际需求增加样式、变量、等文件
```

Options API 示例

```

<templates>
  <div>
    <t-page-header>index.vue示例</t-page-header>
  </div>
</templates>
<script>
export default {
  components: {},
  data() {
    return {};
  },
  methods: {},
};
</script>
<style lang="less">
// 如果需要导入样式
@import "./index.less";

//...
</style>
```

Composition API 示例

```

<templates>
  <div>
    <t-page-header>index.vue示例</t-page-header>
  </div>
</templates>
<script setup>
import { ref, onMounted } from "vue";

// 定义变量
const count = ref(0);

// 定义方法
function increment() {
  count.value++;
}

// 生命周期钩子
onMounted(() => {
  console.log(`The initial count is ${count.value}.`);
});
</script>
<style lang="less">
// 如果需要导入样式
@import "./index.less";

//...
</style>
```

**tips: 一般情况下推荐您使用`Composition API`进行开发，`Composition API`有关的好处请[点击此处](https://vuejs.org/guide/introduction.html#api-styles)**

然后，需要在配置新页面的路由。根据具体的需求，修改 `src/router/modules` 中的文件。

```
export default [
  // ...其他路由
  {
    path: "/new-page",
    title: "新页面侧边栏标题",
    component: "../layouts/td-layout.tsx",
    redirect: "/login/index",
    children: [
      {
        title: "新页面",
        path: "index",
        meta: { needLogin: false },
        component: "../pages/new-page/index.vue",
      },
    ],
  },
];
```

配置后，就可以在项目的侧边栏中找到新页面的入口了。

菜单（侧边栏和面包屑）由路由配置自动生成，根据路由变化可自动匹配，开发者无需手动处理这些逻辑。

**tip：如果您对 vue 的开发方式不是很熟悉，可以查阅 [新手知识](https://vuejs.org/)。**

#### 开发组件[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E5%BC%80%E5%8F%91%E7%BB%84%E4%BB%B6)

当 TDesign 提供的组件不能满足您的需求的时候，您可以根据需要开发新的组件, 推荐放置在`src/component`目录下。

组件的开发方式和 **页面组件** 的开发方式类似，不过您不再需要去为它增加路由，而是在您的组件中引用即可。

首先，在 `src/components` 下新增一个组件文件，`new-component.vue`

```

<template>
  <div>
    <slot name="new-component" />
    <slot />
  </div>
</template>
```

然后，在页面组件中去引用这个组件

Options API 示例

```

<template>
  <div>
    <t-page-header>个人中心</t-page-header>
    
    <my-component v-slot="{ 'new-component':'我插入slot组件的内容' }">
    </my-component>
  </div>
</template>
<script>
// 引入组件
import MyComponent from "@/components/new-component.vue";

export default {
  // 注册组件
  components: {
    MyComponent,
  },
  data() {
    return {};
  },
  methods: {},
};
</script>

<style lang="less">
// 如果需要导入样式
@import "./index.less";

//...
</style>
```

Composition API 示例

```

<template>
  <div>
    <t-page-header>个人中心</t-page-header>
    
    <my-component v-slot="{ 'new-component':'我插入slot组件的内容' }">
    </my-component>
  </div>
</template>
<script setup>
// 引入组件
import MyComponent from "@/components/new-component.vue";
</script>
<style lang="less">
// 如果需要导入样式
@import "./index.less";

//...
</style>
```

这样，一个简单的组件就可以投入使用了。

**tip 如果您对 vue 的开发方式不是很熟悉，可以查阅 [新手知识](https://vuejs.org/)。**

### 布局配置[](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E5%B8%83%E5%B1%80%E9%85%8D%E7%BD%AE)

网站布局支持空布局、侧边栏导航布局、 侧边栏布局加头部导航和头部导航等四种中后台常用布局。布局文件位于 `./src/layouts`。

使用这些布局，您只需要在 `src/router` 中配置路由的时候，将 `父级路由` 配置成相应的布局组件就可以了。示例代码如下：

```
export default [
  {
    path: "/login",
    title: "登录页",
    component: "../layouts/blank.vue", // 这里配置成需要的布局组件
    icon: "chevron-right-rectangle",
    redirect: "/login/index",
    children: [
      {
        title: "登录中心",
        path: "index",
        meta: { needLogin: false },
        component: "../pages/login/index.vue",
      },
    ],
  },
];
```


## 开发规范
为了维护项目的代码质量，项目中内置了格式化代码的工具 `Prettier` 和代码检测质量检查工具 `ESlint`。

同时，也推荐您在开发过程中遵循提交规范，以保持项目仓库的分支、提交信息的清晰整洁。

### 代码编写规范[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#%E4%BB%A3%E7%A0%81%E7%BC%96%E5%86%99%E8%A7%84%E8%8C%83)

#### [Prettier](https://prettier.io/)[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#prettier)

如果您希望项目中的代码都符合统一的格式，推荐您在 VS Code 中安装 `Prettier` 插件。它可以帮助您在每次保存时自动化格式化代码。

在脚手架搭建好的项目中，已经内置了符合 TDesign 开发规范的 `.prettierrc.js` 文件。您只需要安装 `Prettier` 插件，即可使用代码自动格式化的功能。

#### [ESlint](https://eslint.org/)[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#eslint)

`ESlint`可以用来检查代码质量和风格问题。

在脚手架搭建好的项目中，也已经内置了 `.eslintrc` 文件。您可以通过下面命令来进行代码检查和自动修复。

执行以下指令，会进行问题的检查及修复，在 pre-commit 的 git hook 中，项目也内置了检查指令，帮助您在提交代码前发现问题。

```
# 代码检查
npm run lint

# 自动修复问题
npm run lint:fix
```

### 目录的命名规范[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#%E7%9B%AE%E5%BD%95%E7%9A%84%E5%91%BD%E5%90%8D%E8%A7%84%E8%8C%83)

1.目录名全部使用小写，单词需采用复数形式，`kebab-case`形式命名，如果需要有多个单词表达，使用中划线连接。如`new-page`、`components`。

### 文件的命名规范[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#%E6%96%87%E4%BB%B6%E7%9A%84%E5%91%BD%E5%90%8D%E8%A7%84%E8%8C%83)

文件的命名规范按照不同情况进行命名

1.如果该文件是单文件组件/类，采用`PascalCase`形式命名，方便导入和使用。如`TDesignSelect.vue`

2.如果该文件是目录下的主文件，采用 index 名称命名，方便导入。如 `index.ts`, `index.vue`

3.如果该文件是接口定义文件，采用`camelCase`形式命名，方便区分文件关联性。如 `list.ts` 和 `listModel.ts`

4.如果该文件是资源/样式文件，采用`kebab-case`形式命名。

### 类及接口的命名规范[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#%E7%B1%BB%E5%8F%8A%E6%8E%A5%E5%8F%A3%E7%9A%84%E5%91%BD%E5%90%8D%E8%A7%84%E8%8C%83)

1.采用`PascalCase`形式命名。

### 分支规范[](https://tdesign.tencent.com/starter/docs/vue-next/develop-rules#%E5%88%86%E6%94%AF%E8%A7%84%E8%8C%83)

- 主干分支 -- `develop`
- 功能分支 -- `feature`
- 修复分支 -- `hotfix`

`develop`分支只接受通过 Merge Request 合入功能分支。

为保证提交的记录干净整洁，其他分支合入之前需要先 `rebase` develop 分支。

**分支命名规则**：`feature/20210401_功能名称`。


## 路由与菜单

路由与菜单的管理，是前端项目中非常重要的一部分。

为了减少开发配置和理解成本，在 TDesign Starter 项目中，管理菜单路由都规范在`src/router` 这个目录下进行配置。

**tips: 通常情况下不需要去理解和修改`index.ts`, 只需要在`modules`目录下增删文件，即可自动添加更新路由**

配置内容是一个对应菜单层级的可嵌套的数组，如

```
[
  {
    path: "/list",
    name: "list",
    component: Layout,
    redirect: "/list/base",
    meta: { title: "列表页", icon: ListIcon, expanded: true },
    children: [
      {
        path: "base",
        name: "ListBase",
        component: () => import("@/pages/list/base/index.vue"),
        meta: { title: "基础列表页", orderNo: 0 },
      },
      {
        path: "card",
        name: "ListCard",
        component: () => import("@/pages/list/card/index.vue"),
        meta: { title: "卡片列表页", hiddenBreadcrumb: true },
      },
      {
        path: "filter",
        name: "ListFilter",
        component: () => import("@/pages/list/filter/index.vue"),
        meta: { title: "筛选列表页" },
      },
      {
        path: "tree",
        name: "ListTree",
        component: () => import("@/pages/list/tree/index.vue"),
        meta: { title: "树状筛选列表页" },
      },
    ],
  },
];
```

数组中每个配置字段都有具体的作用：

- `path` 是当前路由的路径，会与配置中的父级节点的 path 组成该页面路由的最终路径；如果需要跳转外部链接，可以将`path`设置为 http 协议开头的路径。
- `name` 影响多标签 Tab 页的 keep-alive 的能力，如果要确保页面有 keep-alive 的能力，请保证该路由的`name`与对应页面（SFC)的`name`保持一致。
- `component` 渲染该路由时使用的页面组件
- `redirect` 重定向的路径
- `meta` 主要用途是路由在菜单上展示的效果的配置
    - `meta.title` 该路由在菜单上展示的标题
    - `meta.icon` 该路由在菜单上展示的图标
    - `meta.expanded` 决定该路由在菜单上是否默认展开
    - `meta.orderNo` 该路由在菜单上展示先后顺序，数字越小越靠前，默认为零
    - `meta.hidden` 决定该路由是否在菜单上进行展示
    - `meta.hiddenBreadcrumb` 如果启用了面包屑，决定该路由是否在面包屑上进行展示
    - `meta.single` 如果是多级菜单且只存在一个节点，想在菜单上只展示一级节点，可以使用该配置。_请注意该配置需配置在父节点_
    - `meta.frameSrc` 内嵌 iframe 的地址
    - `meta.frameBlank` 内嵌 iframe 的地址是否以新窗口打开
    - `meta.keepAlive` 可决定路由是否开启keep-alive，默认开启。
- `children` 子菜单的配置

由于 TDesign 菜单的限制，最多只允许配置到`三级菜单`。如果菜单层级超过三级，我们建议梳理业务场景，判断层级是否合理。

由于设计美观需要，官网示例只展示了二级菜单，如果存在三级的配置需求，可以参考以下的代码进行配置：

**tips: 务必注意，三级菜单需要在二级菜单中的组件包含`<router-view>`标签才能正常显示，[详情](https://router.vuejs.org/zh/guide/essentials/nested-routes.html)**

```
{
 path: '/menu',
 name: 'menu',
 component: Layout,
 meta: { title: '一级菜单', icon: 'menu-fold' },
 children: [
    {
      path: 'second',
      meta: { title: '二级菜单' },
      component: () => import('@/layouts/blank.vue'),
      children: [
           {
             path: 'third',
             name: 'NestMenu',
             component: () => import('@/pages/nest-menu/index.vue'),
             meta: { title: '三级菜单' },
           },
      ],
    },
  ],
},
```

## 权限控制

许多系统需要通过权限，控制用户有哪些权限访问部分菜单和路由，常见的控制权限的方式有`后端权限控制`和`前端权限控制`。

### 后端权限控制[](https://tdesign.tencent.com/starter/docs/vue-next/permission-control#%E5%90%8E%E7%AB%AF%E6%9D%83%E9%99%90%E6%8E%A7%E5%88%B6)

在 TDesign Vue Next Starter 0.7.0 版本开始，我们将默认权限控制的方式统一为`后端权限控制`。

通过后端权限控制，可以达到更细颗粒度的权限控制，包括图标、顺序、菜单命名等细节。

使用后端权限控制，需要后端配合一个菜单请求的接口，根据用户身份信息，返回具体的序列化后的菜单列表，模板会将它转换为路由和菜单。 由于是序列化的菜单列表，与[路由与菜单](https://tdesign.tencent.com/starter/docs/vue-next/router-menu)章节相比，需要在返回的菜单接口中将几个非序列化的字段进行序列化。

- `component` 字段：
    
    - 非具体页面路由，默认提供了`LAYOUT`、`BLANK`和`IFRAME`
    - 具体页面路由，请设置为对应页面在项目中的相对路径，如基础列表页对应的是`/list/base/index`
- `meta.icon` 字段：请直接使用 TDesign 的 icon 的中划线命名，如`view-list`，所有图标可以在 [TDesign 图标列表](https://tdesign.tencent.com/vue/components/icon#%E5%85%A8%E9%83%A8%E5%9B%BE%E6%A0%87) 中找到。
    
    **tips:此处图标的序列化是借助了 vite 3+ 的能力引入 node_modules 中的第三方包，会根据 name 引入对应的包内的图标 不会发起网络请求。**
    

序列化后的菜单列表示例如下所示，或可以参考此接口进行返回 👉🏻 [请求菜单列表](https://service-bv448zsw-1257786608.gz.apigw.tencentcs.com/api/get-menu-list)

```
[
  {
    path: "/list",
    name: "list",
    component: "LAYOUT",
    redirect: "/list/base",
    meta: {
      title: "列表页",
      icon: "view-list",
    },
    children: [
      {
        path: "base",
        name: "ListBase",
        component: "/list/base/index",
        meta: {
          title: "基础列表页",
        },
      },
      {
        path: "card",
        name: "ListCard",
        component: "/list/card/index",
        meta: {
          title: "卡片列表页",
        },
      },
      {
        path: "filter",
        name: "ListFilter",
        component: "/list/filter/index",
        meta: {
          title: "筛选列表页",
        },
      },
      {
        path: "tree",
        name: "ListTree",
        component: "/list/tree/index",
        meta: {
          title: "树状筛选列表页",
        },
      },
    ],
  },
];
```

### 前端权限控制[](https://tdesign.tencent.com/starter/docs/vue-next/permission-control#%E5%89%8D%E7%AB%AF%E6%9D%83%E9%99%90%E6%8E%A7%E5%88%B6)

如果您需要使用`前端权限控制`，我们也提供了一个雏形的前端权限控制版本，您可以通过替换`store/permission.ts`的内容为`store/permission-fe.ts`的内容来实现。

在此权限控制下，请将系统可能涉及到的菜单都在`router/modules`下参考固定路由，按项目的具体需求准备好。不需要发起菜单请求，通过用户的 roles 字段中允许访问的菜单，达到对菜单进行过滤筛选，只能访问部分菜单的效果。

## 请求与数据

### 发起请求[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E5%8F%91%E8%B5%B7%E8%AF%B7%E6%B1%82)

TDesign Starter 初始化的项目中，采用 **[axios](https://github.com/axios/axios)** 做为请求的资源库，并对其做了封装，可以从`src/utils/request`的路径中引入封装的 request，并在具体场景中使用。我们建议您在`src/apis`目录中管理您的项目使用到的 api，并在具体组件/页面中使用。 大部分情况下，您不需要改动`src/utils/request`中的代码，只需要在`src/apis`目录中新增您使用的接口，并在页面中引入接口使用即可。

```
// src/apis 管理api请求
import { request } from "@/utils/request";
import type { CardListResult, ListResult } from "@/api/model/listModel";

const Api = {
  BaseList: "/get-list",
  CardList: "/get-card-list",
};

export function getList() {
  return (
    request.get <
    ListResult >
    {
      url: Api.BaseList,
    }
  );
}

export function getCardList() {
  return (
    request.get <
    CardListResult >
    {
      url: Api.CardList,
    }
  );
}
```

```
// src/pages/list/card 引入接口并使用
import { getCardList } from "@/api/list";

const fetchData = async () => {
  try {
    const { list } = await getCardList();
    productList.value = list;
    pagination.value = {
      ...pagination.value,
      total: list.length,
    };
  } catch (e) {
    console.log(e);
  } finally {
    dataLoading.value = false;
  }
};
```

### 请求代理[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E8%AF%B7%E6%B1%82%E4%BB%A3%E7%90%86)

项目中默认启用了直连代理模式，`.env`配置文件的中的`VITE_IS_REQUEST_PROXY`环境变量是对应是否启用直连代理模式的开关，环境变量的具体内容请查看 **[进入开发-环境变量](https://tdesign.tencent.com/starter/docs/vue-next/develop#%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)** 章节。

**tips: 如果`VITE_IS_REQUEST_PROXY`为`true`则采用该配置文件中的地址请求，会绕过`vite.config.js`中设置的代理**

您可以在关闭直连代理模式后，在`vite.config.js`中进行代理设置，使用 **Vite** 的`http-proxy`。

- 示例：

```
export default defineConfig({
  server: {
    proxy: {
      // 字符串简写写法
      "/foo": "http://localhost:4567/foo",
      // 选项写法
      "/api": {
        target: "http://jsonplaceholder.typicode.com",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      // 正则表达式写法
      "^/fallback/.*": {
        target: "http://jsonplaceholder.typicode.com",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/fallback/, ""),
      },
      // 使用 proxy 实例
      "/api": {
        target: "http://jsonplaceholder.typicode.com",
        changeOrigin: true,
        configure: (proxy, options) => {
          // proxy 是 'http-proxy' 的实例
        },
      },
    },
  },
});
```

完整选项详见 [http-party 的配置](https://github.com/http-party/node-http-proxy#options)。

### Mock 数据[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#mock-%E6%95%B0%E6%8D%AE)

如果需要进行数据 Mock，在 `vite.config.js` 中，将 viteMockServe 中配置 `localEnabled` 为 `true` ，即可开启 mock server 的拦截。

```
viteMockServe({
    mockPath: 'mock',
    localEnabled: true,
}),
```

### 高级配置-部分请求不代理的场景[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AE-%E9%83%A8%E5%88%86%E8%AF%B7%E6%B1%82%E4%B8%8D%E4%BB%A3%E7%90%86%E7%9A%84%E5%9C%BA%E6%99%AF)

在某些业务场景下可能会使用到腾讯云的 COS 对象存储或其他厂商的上传服务，在此情况下则无法直接使用`@/utils/request`进行请求，否则地址会被代理。

此情况下可以在`src/utils/request/index.ts`中最下方添加新的请求实例

- 示例：

```
function createOtherAxios(opt?: Partial<CreateAxiosOptions>) {
  return new VAxios(
    merge(
      <CreateAxiosOptions>{
        // https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#authentication_schemes
        // 例如: authenticationScheme: 'Bearer'
        authenticationScheme: '',
        // 超时
        timeout: 10 * 1000,
        // 携带Cookie
        withCredentials: true,
        // 头信息
        headers: { 'Content-Type': 'application/json;charset=UTF-8' },
        // 数据处理方式
        transform,
        // 配置项，下面的选项都可以在独立的接口请求中覆盖
        requestOptions: {
          // 接口地址
          apiUrl: '',
          // 是否自动添加接口前缀
          isJoinPrefix: false,
          // 接口前缀
          // 例如: https://www.baidu.com/api
          // urlPrefix: '/api'
          urlPrefix: '',
          // 是否返回原生响应头 比如：需要获取响应头时使用该属性
          isReturnNativeResponse: false,
          // 需要对返回数据进行处理
          isTransformResponse: false,
          // post请求的时候添加参数到url
          joinParamsToUrl: false,
          // 格式化提交参数时间
          formatDate: true,
          // 是否加入时间戳
          joinTime: true,
          // 忽略重复请求
          ignoreRepeatRequest: true,
          // 是否携带token
          withToken: true,
          // 重试
          retry: {
            count: 3,
            delay: 1000,
          },
        },
      },
      opt || {},
    ),
  );
}
export const requestOther = createOtherAxios();
```

在添加新实例后，引入新实例`@/utils/requestOther`即可继续开发

### 高级配置-不需要重试的场景[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AE-%E4%B8%8D%E9%9C%80%E8%A6%81%E9%87%8D%E8%AF%95%E7%9A%84%E5%9C%BA%E6%99%AF)

此情况下可以在`src/utils/request/index.ts`中最下方的`createAxios`方法中的参数`retry`移除即可

### 高级配置-修改请求返回的通用模型[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AE-%E4%BF%AE%E6%94%B9%E8%AF%B7%E6%B1%82%E8%BF%94%E5%9B%9E%E7%9A%84%E9%80%9A%E7%94%A8%E6%A8%A1%E5%9E%8B)

首先需要您在`src/types/axios.d.ts`中的`Result`中声明您的通用模型

- 示例：

```
export interface Result<T = any> {
  code: number;
  data: T;
}
```

随后在`src/utils/request/index.ts`中的`transform`方法中对您的数据进行预处理

**tips: 如果您不需要对数据预处理则可以在最下方将`isTransformResponse`设置关闭**

### 高级配置-修改请求params参数的序列化方式[](https://tdesign.tencent.com/starter/docs/vue-next/request-data#%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AE-%E4%BF%AE%E6%94%B9%E8%AF%B7%E6%B1%82params%E5%8F%82%E6%95%B0%E7%9A%84%E5%BA%8F%E5%88%97%E5%8C%96%E6%96%B9%E5%BC%8F)

使用[qs](https://github.com/ljharb/qs)序列化请求params参数

首先需要您在`src/utils/request/Axios.ts`中的`supportParamsStringify`方法中选择您需要的序列化方式

- 示例：

```
// 支持params数组参数格式化
  supportParamsStringify(config: AxiosRequestConfig) {
    const headers = config.headers || this.options.headers;
    const contentType = headers?.['Content-Type'] || headers?.['content-type'];

    if (contentType === ContentTypeEnum.FormURLEncoded || !Reflect.has(config, 'params')) {
      return config;
    }

    return {
      ...config,
      //修改此处的arrayFormat，选项有'indices' 'brackets' 'repeat' 'comma'等，请参考qs文档根据项目需要选择
      paramsSerializer: (params: any) => stringify(params, { arrayFormat: 'brackets' }), 
    };
  }
```

随后在同一文件中的`request`方法中，取消调用`supportParamsStringify`行的注释

```
conf = this.supportParamsStringify(conf);
```

**tips: axios会使用内置的toFormData以brackets方式序列化params参数，`如果您不需要修改，无需进行上述操作`**


## 样式与静态资源

### 本地静态资源存放[](https://tdesign.tencent.com/starter/docs/vue-next/style#%E6%9C%AC%E5%9C%B0%E9%9D%99%E6%80%81%E8%B5%84%E6%BA%90%E5%AD%98%E6%94%BE)

静态资源可以放在 `./src/assets` 目录下，然后在文件中通过相对路径引入。

### 如何引入字体、图片[](https://tdesign.tencent.com/starter/docs/vue-next/style#%E5%A6%82%E4%BD%95%E5%BC%95%E5%85%A5%E5%AD%97%E4%BD%93%E3%80%81%E5%9B%BE%E7%89%87)

#### 字体[](https://tdesign.tencent.com/starter/docs/vue-next/style#%E5%AD%97%E4%BD%93)

将字体文件`.ttf`放在`./src/assets/fonts` 目录下。

然后在文件`./src/style/font-family.less`中引入该字体文件。

```
@font-face {
  font-family: "w7";
  src: url("w7.ttf");
  font-weight: normal;
  font-style: normal;
}
```

在 App.vue 中的 style 里引入

```
<style lang="less" rel="stylesheet/less">
  @import "./src/style/font-family.less";
</style>
```

#### 图片[](https://tdesign.tencent.com/starter/docs/vue-next/style#%E5%9B%BE%E7%89%87)

将图片文件放在 `./src/assets/images` 目录下。

在 vue 文件中通过相对路径引入`@/assets/images/image.png`。

### 引入 SVG[](https://tdesign.tencent.com/starter/docs/vue-next/style#%E5%BC%95%E5%85%A5-svg)

SVG 是一种可变向量图，提供了 DOM 编程的接口，更多关于 SVG 的[点击这里](https://developer.mozilla.org/zh-CN/docs/Web/SVG)

通过源码引入，如下:

```
<template>
  <svg
    width="34"
    height="34"
    viewBox="0 0 34 34"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="0.5" y="0.5" width="33" height="33" rx="16.5" stroke="white" />
    <path
      d="M16.35 17.6501V21.5H17.65V17.6501H21.5V16.3501H17.65V12.5H16.35V16.3501H12.5V17.6501H16.35Z"
    />
  </svg>
</template>
```

通过路径引入，可以像组件一样使用（此功能实现借助了插件 `vite-plugin-vue2-svg` ）

```
<template>
  <t-logow class="t-logo" />
</template>

<script>
import tLogow from "../assets/t-logo-colorful.svg";
export default {
  components: {
    tLogow,
  },
};
</script>
```