import {NgModule} from '@angular/core';
import {NZ_ICONS, NzIconModule} from 'ng-zorro-antd/icon';

import {
  MenuFoldOutline,
  MenuUnfoldOutline,
  FormOutline,
  DashboardOutline,
  SettingOutline,
  LogoutOutline,
  MailOutline
} from '@ant-design/icons-angular/icons';
import {CommonModule} from "@angular/common";

const icons = [MenuFoldOutline, MenuUnfoldOutline, DashboardOutline, FormOutline, SettingOutline, LogoutOutline, MailOutline];

@NgModule({
  imports: [CommonModule, NzIconModule.forChild(icons)],
  exports: [NzIconModule],
  providers: [
    //{ provide: NZ_ICONS, useValue: icons }
  ]
})
export class IconsProviderModule {
}
