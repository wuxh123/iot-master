import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ProjectComponent} from './project/project.component';
import {ProjectAdapterComponent} from './project-adapter/project-adapter.component';
import {ProjectVariableComponent} from './project-variable/project-variable.component';
import {ProjectBatchComponent} from './project-batch/project-batch.component';
import {ProjectJobComponent} from './project-job/project-job.component';
import {ProjectStrategyComponent} from './project-strategy/project-strategy.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {ProjectStrategyEditComponent} from './project-strategy-edit/project-strategy-edit.component';
import {ProjectJobEditComponent} from './project-job-edit/project-job-edit.component';
import {ProjectBatchEditComponent} from './project-batch-edit/project-batch-edit.component';
import {ProjectVariableEditComponent} from './project-variable-edit/project-variable-edit.component';
import {ProjectAdapterEditComponent} from './project-adapter-edit/project-adapter-edit.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ProjectDetailComponent} from './project-detail/project-detail.component';
import {DeviceComponent} from './device/device.component';
import {DeviceDetailComponent} from './device-detail/device-detail.component';
import {DeviceLogComponent} from './device-log/device-log.component';
import {DeviceMapComponent} from './device-map/device-map.component';
import {ElementComponent} from './element/element.component';
import {ElementEditComponent} from './element-edit/element-edit.component';
import {AdapterComponent} from './adapter/adapter.component';
import {AdapterEditComponent} from './adapter-edit/adapter-edit.component';
import {HistoryComponent} from './history/history.component';
import {AlgorithmComponent} from './algorithm/algorithm.component';
import {AlertComponent} from "./alert/alert.component";

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

      {path: 'plugin', component: PluginComponent},
      {path: 'plugin/create', component: PluginEditComponent},
      {path: 'plugin/:id/edit', component: PluginEditComponent},

      {path: 'project', component: ProjectComponent},
      {path: 'project/create', component: ProjectEditComponent},
      {path: 'project/:id/edit', component: ProjectEditComponent},
      {path: 'project/:id/detail', component: ProjectDetailComponent},
      {path: 'project/:id/adapter', component: ProjectAdapterComponent},
      {path: 'project/:id/adapter/create', component: ProjectAdapterEditComponent},
      {path: 'project/:id/adapter/:id/edit', component: ProjectAdapterEditComponent},
      {path: 'project/:id/variable', component: ProjectVariableComponent},
      {path: 'project/:id/variable/create', component: ProjectVariableEditComponent},
      {path: 'project/:id/variable/:id/edit', component: ProjectVariableEditComponent},
      {path: 'project/:id/batch', component: ProjectBatchComponent},
      {path: 'project/:id/batch/create', component: ProjectBatchEditComponent},
      {path: 'project/:id/batch/:id/edit', component: ProjectBatchEditComponent},
      {path: 'project/:id/job', component: ProjectJobComponent},
      {path: 'project/:id/job/create', component: ProjectJobEditComponent},
      {path: 'project/:id/job/:id/edit', component: ProjectJobEditComponent},
      {path: 'project/:id/strategy', component: ProjectStrategyComponent},
      {path: 'project/:id/strategy/create', component: ProjectStrategyEditComponent},
      {path: 'project/:id/strategy/:id/edit', component: ProjectStrategyEditComponent},

      {path: 'element', component: ElementComponent},
      {path: 'element/create', component: ElementEditComponent},
      {path: 'element/:id/edit', component: ElementEditComponent},

      {path: 'adapter', component: AdapterComponent},
      {path: 'adapter/create', component: AdapterEditComponent},
      {path: 'adapter/:id/edit', component: AdapterEditComponent},

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
