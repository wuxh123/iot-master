import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {PasswordComponent} from './password/password.component';
import {LinkComponent} from './link/link.component';
import {ChannelComponent} from './channel/channel.component';
import {UserComponent} from './user/user.component';
import {PluginComponent} from './plugin/plugin.component';
import {PluginStoreComponent} from './plugin-store/plugin-store.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ChannelMonitorComponent} from './channel-monitor/channel-monitor.component';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'channel', component: ChannelComponent},
      {path: 'channel-monitor/:id', component: ChannelMonitorComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link-monitor/:id', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginComponent},
      {path: 'plugin-store', component: PluginStoreComponent},
      {path: 'user', component: UserComponent},
      // {path: 'password', component: PasswordComponent},
      {path: '**', redirectTo: 'dash'},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {
}
