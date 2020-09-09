import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {ChannelComponent} from './channel/channel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'channel', component: ChannelComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link-monitor/:id', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginComponent},
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
