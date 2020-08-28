import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from "./dash/dash.component";
import {MainComponent} from "./main.component";
import {CopyComponent} from "./copy/copy.component";
import {SettingComponent} from "./setting/setting.component";
import {AuditComponent} from "./audit/audit.component";
import {PasswordComponent} from "./password/password.component";
import {UploadComponent} from "./upload/upload.component";
import {EmailComponent} from "./email/email.component";
import {SubscribeComponent} from "./subscribe/subscribe.component";
import {DownloadComponent} from "./download/download.component";
import {KeywordComponent} from "./keyword/keyword.component";

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'copy', component: CopyComponent},
      {path: 'upload', component: UploadComponent},
      {path: 'download', component: DownloadComponent},
      {path: 'subscribe', component: SubscribeComponent},
      {path: 'email', component: EmailComponent},
      {path: 'audit', component: AuditComponent},
      {path: 'keyword', component: KeywordComponent},
      {path: 'setting', component: SettingComponent},
      {path: 'password', component: PasswordComponent},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {
}
