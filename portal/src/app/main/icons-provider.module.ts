import {NgModule} from '@angular/core';
import {NZ_ICONS, NzIconModule} from 'ng-zorro-antd/icon';

import {
  MenuFoldOutline,
  MenuUnfoldOutline,
  DashboardOutline,
  SettingOutline,
  LogoutOutline,
  ApiOutline,
  ReloadOutline,
  PlusOutline,
  DeleteOutline,
  AppstoreOutline,
  AimOutline,
  SwapOutline, ClusterOutline,
} from '@ant-design/icons-angular/icons';
import {CommonModule} from '@angular/common';

const icons = [
  // 菜单相关
  MenuFoldOutline, MenuUnfoldOutline, DashboardOutline, ApiOutline, SettingOutline, AppstoreOutline,
  // 表格操作
  ReloadOutline, PlusOutline, DeleteOutline, AimOutline, SwapOutline,
  LogoutOutline, ClusterOutline
];

@NgModule({
  imports: [CommonModule, NzIconModule.forChild(icons)],
  exports: [NzIconModule],
  providers: [
    {provide: NZ_ICONS, useValue: icons}
  ]
})
export class IconsProviderModule {
}
