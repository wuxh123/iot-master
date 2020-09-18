import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-plugin-edit',
  templateUrl: './plugin-edit.component.html',
  styleUrls: ['./plugin-edit.component.scss']
})
export class PluginEditComponent implements OnInit {

  @Input() plugin: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.plugin.id) {
      this.as.get('plugin/' + this.plugin.id).subscribe(res => {
        this.plugin = res.data;
      });
    }
  }

  submit(): void {
    if (this.plugin.id) {
      this.as.put('plugin/' + this.plugin.id, this.plugin).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('plugin', this.plugin).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
