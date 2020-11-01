import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkComponent} from './link/link.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {PluginComponent} from './plugin/plugin.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ProjectComponent} from './project/project.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {ProjectDetailComponent} from './project-detail/project-detail.component';
import {ElementComponent} from './element/element.component';
import {ElementDetailComponent} from './element-detail/element-detail.component';
import {ElementEditComponent} from './element-edit/element-edit.component';
import {HistoryComponent} from './history/history.component';
import {AlgorithmComponent} from './algorithm/algorithm.component';
import {AlertComponent} from './alert/alert.component';
import {DeviceComponent} from './device/device.component';
import {DeviceMapComponent} from './device-map/device-map.component';
import {DeviceLogComponent} from './device-log/device-log.component';
import {DeviceDetailComponent} from './device-detail/device-detail.component';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},

      {path: 'device', component: DeviceComponent},
      {path: 'device/map', component: DeviceMapComponent},
      {path: 'device/log', component: DeviceLogComponent},
      {path: 'device/:id/detail', component: DeviceDetailComponent},
      {path: 'device/:id/log', component: DeviceLogComponent},

      {path: 'history', component: HistoryComponent},
      {path: 'algorithm', component: AlgorithmComponent},

      {path: 'alert', component: AlertComponent},

      {path: 'tunnel', component: TunnelComponent},
      {path: 'tunnel/create', component: TunnelEditComponent},
      {path: 'tunnel/:id/edit', component: TunnelEditComponent},
      {path: 'tunnel/:id/link', component: LinkComponent},

      {path: 'link', component: LinkComponent},
      {path: 'link/:id/edit', component: LinkEditComponent},
      {path: 'link/:id/monitor', component: LinkMonitorComponent},

      {path: 'project', component: ProjectComponent},
      {path: 'project/create', component: ProjectEditComponent},
      {path: 'project/:id/edit', component: ProjectEditComponent},
      {path: 'project/:id/detail', component: ProjectDetailComponent},

      {path: 'element', component: ElementComponent},
      {path: 'element/create', component: ElementEditComponent},
      {path: 'element/:id/edit', component: ElementEditComponent},
      {path: 'element/:id/detail', component: ElementDetailComponent},


      {path: 'plugin', component: PluginComponent},
      {path: 'plugin/create', component: PluginEditComponent},
      {path: 'plugin/:id/edit', component: PluginEditComponent},

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
