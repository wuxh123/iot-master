import {NgModule} from '@angular/core';

import {IconsProviderModule} from './icons-provider.module';
import {NzLayoutModule} from 'ng-zorro-antd/layout';
import {NzMenuModule} from 'ng-zorro-antd/menu';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {MainRoutingModule} from './main-routing.module';

import {MainComponent} from './main.component';
import {
  NzButtonModule,
  NzCheckboxModule, NzDividerModule, NzDrawerModule,
  NzFormModule,
  NzIconModule,
  NzInputModule, NzInputNumberModule,
  NzModalModule, NzPopconfirmModule, NzSelectModule, NzSwitchModule,
  NzTableModule,
  NzToolTipModule
} from 'ng-zorro-antd';
import {DashComponent} from './dash/dash.component';
import {MomentModule} from 'ngx-moment';
import {UiModule} from '../ui/ui.module';
import {ChannelComponent} from './channel/channel.component';
import {LinkComponent} from './link/link.component';
import {UserComponent} from './user/user.component';
import {PasswordComponent} from './password/password.component';
import {ChannelEditComponent} from './channel-edit/channel-edit.component';
import {ChannelDetailComponent} from './channel-detail/channel-detail.component';
import {LinkDetailComponent} from './link-detail/link-detail.component';
import {UserEditComponent} from './user-edit/user-edit.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {NzSpaceModule} from "ng-zorro-antd/space";
import {LinkEditComponent} from "./link-edit/link-edit.component";


@NgModule({
  declarations: [MainComponent, DashComponent,
    ChannelComponent, ChannelDetailComponent, ChannelEditComponent,
    LinkComponent, LinkDetailComponent, LinkMonitorComponent, LinkEditComponent,
    UserComponent, UserEditComponent,
    PasswordComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MomentModule,
    UiModule,
    // IconsProviderModule,
    // BrowserModule,
    NzIconModule,
    NzLayoutModule,
    NzMenuModule,
    HttpClientModule,
    MainRoutingModule,
    NzToolTipModule,
    NzTableModule,
    NzModalModule,
    NzFormModule,
    NzButtonModule,
    NzInputModule,
    NzCheckboxModule,
    NzSwitchModule,
    NzPopconfirmModule,
    IconsProviderModule,
    NzDividerModule,
    NzDrawerModule,
    NzSelectModule,
    NzSpaceModule,
    NzInputNumberModule,
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
