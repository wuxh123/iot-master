import {NgModule} from '@angular/core';
import {NZ_ICONS, NzIconModule} from 'ng-zorro-antd/icon';

import {
  MenuFoldOutline,
  MenuUnfoldOutline,
  FormOutline,
  DashboardOutline,
  SettingOutline,
  LogoutOutline,
  ApiOutline,
  HddOutline,
  ApartmentOutline,
  ReloadOutline,
  PlusOutline,
  DeleteOutline,
} from '@ant-design/icons-angular/icons';
import {CommonModule} from '@angular/common';

const icons = [MenuFoldOutline, MenuUnfoldOutline, DashboardOutline,
  FormOutline, SettingOutline, LogoutOutline, ApiOutline, ReloadOutline, PlusOutline, DeleteOutline,
  HddOutline, ApartmentOutline];

@NgModule({
  imports: [CommonModule, NzIconModule.forChild(icons)],
  exports: [NzIconModule],
  providers: [
    {provide: NZ_ICONS, useValue: icons}
  ]
})
export class IconsProviderModule {
}
