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
  NzCheckboxModule,
  NzFormModule,
  NzIconModule,
  NzInputModule,
  NzModalModule, NzPopconfirmModule, NzSwitchModule,
  NzTableModule,
  NzToolTipModule
} from 'ng-zorro-antd';
import {CopyComponent} from './copy/copy.component';
import {DashComponent} from './dash/dash.component';
import {UploadComponent} from './upload/upload.component';
import {DownloadComponent} from './download/download.component';
import {KeywordComponent} from './keyword/keyword.component';
import {AuditComponent} from './audit/audit.component';
import {SubscribeComponent} from './subscribe/subscribe.component';
import {EmailComponent} from './email/email.component';
import {SettingComponent} from './setting/setting.component';
import {PasswordComponent} from './password/password.component';
import {MomentModule} from 'ngx-moment';
import {UiModule} from '../ui/ui.module';
import {DetailComponent} from './detail/detail.component';


@NgModule({
  declarations: [MainComponent, DashComponent, CopyComponent, DetailComponent,
    UploadComponent, DownloadComponent, KeywordComponent,
    AuditComponent, SubscribeComponent, EmailComponent,
    SettingComponent, PasswordComponent],
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
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
